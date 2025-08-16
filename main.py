import os
import logging
import time
from dotenv import load_dotenv # pip install python-dotenv playwright / playwright install 필요
from playwright.sync_api import sync_playwright, TimeoutError as PlaywrightTimeoutError
import requests

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

def load_environment():
    """
    Load environment variables from .env file.
    """
    load_dotenv()
    discord_url = os.getenv("DISCORD_URL")
    if not discord_url:
        logging.error("DISCORD_URL not found in environment variables.")
        raise EnvironmentError("DISCORD_URL not found in environment variables.")
    return discord_url

def send_discord_message(webhook_url, content):
    """
    Send a message to Discord using the provided webhook URL.
    """
    data = {
        "content": content,
        # "username": "Event Bot",
        # "avatar_url": "https://i.imgur.com/4M34hi2.png"  # Optional: Change to your desired avatar
    }
    try:
        response = requests.post(webhook_url, json=data)
        if response.status_code != 204:
            logging.error(f"Failed to send Discord message: {response.status_code} - {response.text}")
    except Exception as e:
        logging.error(f"Exception while sending Discord message: {e}")

def extract_events(page):
    """
    Extract event data from the webpage.
    """
    events = []
    try:
        # Navigate to the events page
        web_url = "https://dev-event.vercel.app/events"
        logging.info(f"Navigating to {web_url}")
        page.goto(web_url, timeout=60000)  # 60 seconds timeout
        page.wait_for_selector("[class^='Home_section__']", timeout=30000)  # .Home_section__EaDnq / 페이지 로딩 대기 30 seconds timeout

        # Wait additional time to ensure all events are loaded
        time.sleep(3)

        # Get all event containers
        event_nodes = page.query_selector_all("[class^='Item_item__container']") # .Item_item__container___T09W
        logging.info(f"Found {len(event_nodes)} event nodes.")

        for node in event_nodes:
            try:
                # Extract the D-Day text
                dday_element = node.query_selector("[class^='DdayTag_tag__']") # .DdayTag_tag__6_oE7
                if not dday_element:
                    logging.debug("D-Day element not found, skipping event.")
                    continue
                today_text = dday_element.text_content().strip()
                if "Today" not in today_text:
                    continue  # Skip events not happening today

                # Extract all tags and check for 'AI'
                tag_elements = node.query_selector_all("[class^='FilterTag_tag__']") # .FilterTag_tag__etNfv
                is_ai = any("AI" in tag.text_content() for tag in tag_elements)
                if not is_ai:
                    continue  # Skip events without 'AI' tag

                # Extract title
                title_element = node.query_selector("[class^='Item_item__content__title__']") # .Item_item__content__title__94_8Q
                title = title_element.text_content().strip() if title_element else "No Title"

                # Extract URL
                link_element = node.query_selector("a")
                url = "<https://dev-event.vercel.app" + link_element.get_attribute("href").strip() + ">" if link_element else "No URL"

                # Extract date
                date_element = node.query_selector("[class^='Item_date__date__']") # .Item_date__date__CoMqV
                date_text = date_element.text_content().strip() if date_element else "No Date"

                # Extract host
                host_element = node.query_selector("[class^='Item_host__']") # .Item_host__3dy8_
                host_text = host_element.text_content().strip() if host_element else "No Host"

                # Append the event data
                event = {
                    "title": title,
                    "desc": f"{url}\n주최: {host_text}\n모집: {date_text}"
                }
                events.append(event)
                logging.debug(f"Extracted event: {event}")

            except Exception as e:
                logging.warning(f"Error extracting event data: {e}")
                continue

    except PlaywrightTimeoutError: # 페이지 로딩 에러
        logging.error("Timeout while loading the events page.")
    except Exception as e: # 그 외 에러
        logging.error(f"Unexpected error: {e}")

    logging.info(f"Total filtered events: {len(events)}")
    return events

def build_discord_message(events):
    """
    Build the Discord message content from the list of events.
    """
    if not events:
        return "## 진행중인 개발자 행사가 없습니다."

    content = "## 진행중인 개발자 행사\n"
    for event in events:
        content += f"### {event['title']}\n"
        content += f"{event['desc']}\n\n"
    return content

def main():
    # Load environment variables (= Discord URL)
    try:
        discord_webhook_url = load_environment()
    except EnvironmentError as e:
        logging.critical(e)
        return

    # Initialize Playwright and extract events
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        context = browser.new_context()
        page = context.new_page()
        events = extract_events(page)
        browser.close()

    # Build the message content
    message_content = build_discord_message(events)
    logging.info("Constructed Discord message.")
    # print(message_content)
    # Send the message to Discord
    send_discord_message(discord_webhook_url, message_content)
    logging.info("Discord message sent successfully.")

if __name__ == "__main__":
    main()