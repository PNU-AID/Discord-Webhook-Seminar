# AI 행사 알림 봇 (Discord Webhook)

이 스크립트는 [DevEvent](https://dev-event.vercel.app/events) 페이지에서 AI 관련 행사를 스크래핑하여 Discord Webhook을 통해 특정 Discord 채널로 알림을 전송합니다.

---

## 기능

- DevEvent 페이지에서 **오늘** 진행 중인 AI 관련 행사를 스크래핑.
- **`AI` 태그**가 포함된 행사만 필터링.
- 행사 제목, 링크, 주최, 모집 기간 등의 정보를 Discord로 전송.
- 환경 변수로 간단히 Discord URL 설정 가능.

## 출력 예시
### 행사가 있는 경우
```
## 진행중인 개발자 행사
### AI 워크샵 2023
https://example.com/event-link
주최: Dev Community
모집: 2023-12-15 ~ 2023-12-18

### AI 컨퍼런스
https://example.com/event-link2
주최: Tech Society
모집: 2023-12-15
```
### 행사가 없는 경우
```
## 진행중인 개발자 행사가 없습니다.
```

## 요구 사항
* Python 3.7 이상
* playwright
* python-dotenv
* requests
