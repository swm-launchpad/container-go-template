# Test Projects for Container Templates

이 디렉토리는 각 템플릿의 동작을 검증하기 위한 최소한의 Hello World 테스트 프로젝트를 포함합니다.

## 구조

```
test-projects/
├── frontend/           # 프론트엔드 프레임워크 (4개)
│   ├── vuejs/         # Vue.js 3 + Vite
│   ├── react/         # React 18 + Vite
│   ├── nextjs/        # Next.js 14 App Router
│   └── static-html/   # Static HTML
├── backend/           # 백엔드 프레임워크 (7개)
│   ├── expressjs/     # Express.js 4
│   ├── nestjs/        # NestJS 10
│   ├── spring-boot/   # Spring Boot 3.2 (Java)
│   ├── kotlin-spring-boot/ # Spring Boot 3.2 (Kotlin)
│   ├── fastapi/       # FastAPI
│   ├── flask/         # Flask 3
│   └── django/        # Django 5
└── database/          # 데이터베이스 (4개)
    ├── mysql/         # MySQL 8
    ├── postgresql/    # PostgreSQL 16
    ├── mongodb/       # MongoDB 7
    └── redis/         # Redis 7
```

## 용도

1. **템플릿 검증**: 각 템플릿이 올바르게 동작하는지 확인
2. **문서화 참고**: 사용자에게 제공할 예제 코드
3. **CI/CD 테스트**: 자동화된 템플릿 테스트에 활용 가능

## 테스트 방법

각 프로젝트 디렉토리의 README.md 또는 소스 코드를 참고하여 실행합니다.

### 프론트엔드 프로젝트
```bash
cd frontend/vuejs
npm install
npm run dev
```

### 백엔드 프로젝트 (Node.js)
```bash
cd backend/expressjs
npm install
npm start
```

### 백엔드 프로젝트 (Python)
```bash
cd backend/fastapi
pip install -r requirements.txt
uvicorn main:app --reload
```

### 백엔드 프로젝트 (Java/Kotlin)
```bash
cd backend/spring-boot
mvn spring-boot:run
```

### 데이터베이스 프로젝트
각 데이터베이스는 Docker 컨테이너로 실행하며, init 스크립트가 포함되어 있습니다.

## 유지보수

- 템플릿 버전 업데이트 시 해당 테스트 프로젝트도 함께 업데이트
- 각 프로젝트는 최소한의 의존성만 포함
- "Hello World" 수준의 간단한 코드 유지
