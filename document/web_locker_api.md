# 라커-서버간 API 명세서
휴대폰 보관 하드웨어(이하 라커)와 웹 서버(이하 서버)간 통신을 위한 API의 명세서입니다.

## 규칙
- HTTP 프로토콜을 통해 통신합니다.
- REST 아키텍처를 준용해 설계합니다.
- HMAC 인증 내용을 담은 헤더와 요청 본문이 존재할 시 timestamp를 같이 전송해야 합니다.
- API는 라커가 시간을 유지하기 위하거나 가져오기 위한 장치가 없으며, 서버로부터 수신받은 timestamp를 runtime 카운터로 증가시키며 시간을 유지한다고 기장힙니다.

## API
### 보안
요청의 발신자가 불법적이지 않은 유효한 하드웨어임을 인증하기위해 HMAC 알고리즘을 사용합니다.
모든 요청에는 `Auth` 헤더와, 요청 본문의 JSON 루트에 `timestamp` 요소, 그에 맞는 올바른 값을 포함하여 전송되야 합니다. 유효하지 않은 값을 포함했을 경우 요청은 거부됩니다. 
`Auth` 헤더에는 생성된 HMAC 해쉬 값이 포함되어야하며 `timestamp` 요소에는 서버 시간을 기준으로한 Unix timestamp가 포함되어야 합니다. 

**HMAC**의 해쉬 알고리즘은 SHA256를(아두이노에서 [지원](https://rweather.github.io/arduinolibs/classSHA256.html)) 사용합니다. 
해쉬는 비밀키와 HTTP 요청 본문을 결합한 내용으로 생성되며 의사코드로 나타내면 다음과 같습니다.
```
hash = sha256(secure_key + http_body)
```

#### GET /api/timestamp
서버의 Unix timestamp를 얻습니다.

|요구 헤더 이름|헤더 내용|
|--|---|
|없음||

|요구 본문 요소|요소 내용|
|--|---|
|없음||

|응답 코드|응답 본문|설명|
|--|---|---|
|200|없음|정상|

### 라커 정보
#### GET /api/locker/{라커 UID}/tag
라커에 할당된 태그 목록을 얻습니다.

|요구 헤더 이름|헤더 내용|
|--|---|
|Auth|HMAC 해시|

|요구 본문 요소|요소 내용|
|--|---|
|timestamp|서버 시간을 기준으로한 Unix timestamp|

|응답 코드|응답 본문|설명|
|--|---|---|
|200|없음|정상|
|400|없음|요청 본문의 형식이 올바르지 않거나 누락된 요소가 있습니다.|
|403|없음|인증 정보가 유효하지 않습니다.|
|500|없음|서버 내부 오류|

### 반납 정보
#### POST /api/locker/{라커 UID}/tag
라커에 반납된 태그 목록을 서버로 보냅니다.

|요구 헤더 이름|헤더 내용|
|--|---|
|Auth|HMAC 해시|

|요구 본문 요소|요소 내용|
|--|---|
|timestamp|서버 시간을 기준으로한 Unix timestamp|
|tag_uids|반납된 휴대폰의 태그 UID 문자열의 배열입니다.|
|weight|반납된 휴대폰의 무게 총합입니다. 그램 단위의 부동소수점입니다.|

|응답 코드|응답 본문|설명|
|--|---|---|
|200|없음|정상|
|400|없음|요청 본문의 형식이 올바르지 않거나 누락된 요소가 있습니다.|
|403|없음|인증 정보가 유효하지 않습니다.|
|500|없음|서버 내부 오류|

#### POST /api/locker/{라커 UID}/door
라커의 문 개폐 시간을 보냅니다.

```
개방            ----------
(센서 값)       |         |
폐쇄       -----          -----

               ^         ^
             Topen     Tclose
               ~~~~Td~~~~~
```

|요구 헤더 이름|헤더 내용|
|--|---|
|Auth|HMAC 해시|

|요구 본문 요소|요소 내용|
|--|---|
|timestamp|서버 시간을 기준으로한 Unix timestamp|
|close_time|Tclose의 Unix timestamp입니다. 0을 보낼 시 서버는 요청 수신 시간으로 대체합니다.|
|duration|밀리세컨드 단위의 Td 간격입니다. Td가 설정한 시간 이상이면 라커는 이를 감지하고 -1을 보냅니다. 서버에서는 이를 이상 행위로 판단합니다.|

|응답 코드|응답 본문|설명|
|--|---|---|
|200|없음|정상|
|400|없음|요청 본문의 형식이 올바르지 않거나 누락된 요소가 있습니다.|
|403|없음|인증 정보가 유효하지 않습니다.|
|500|없음|서버 내부 오류|
