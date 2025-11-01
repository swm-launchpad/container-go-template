# Container Go Template

Container Dockerfile 템플릿 관리 레포지토리

## 📋 개요

이 레포지토리는 LaunchPad 플랫폼에서 사용하는 Container Dockerfile 템플릿을 관리합니다.

템플릿 변경 시 자동으로 운영 데이터베이스에 마이그레이션을 적용하며, Git diff 기반으로 변경된 템플릿만 업데이트합니다.

## 🏗️ 구조

```
container-go-template/
├── templates/              # 템플릿 디렉토리
│   ├── vuejs/
│   │   ├── template.tmpl   # Dockerfile 템플릿 (Go template 문법)
│   │   ├── config.json     # 템플릿 옵션, 포트, 환경변수
│   │   └── metadata.json   # 이름, 설명, 카테고리, 버전
│   ├── react/
│   ├── nextjs/
│   └── ...                 # 23개 프레임워크 템플릿
├── scripts/
│   └── generate-migration.go  # 마이그레이션 SQL 생성 도구
└── .github/workflows/
    ├── ci.yml              # PR 검증
    └── cd.yml              # 운영 배포
```

## 🎯 현재 지원 템플릿

### Frontend (6)
- Vue.js
- React
- Next.js
- Angular
- Svelte
- Static HTML

### Backend (13)
- Express.js
- NestJS
- Go Gin
- Spring Boot
- Kotlin Spring Boot
- FastAPI
- Flask
- Django
- Laravel
- Ruby on Rails
- Rust
- .NET Core
- Fiber

### Database (4)
- MySQL
- PostgreSQL
- MongoDB
- Redis

## 🔧 템플릿 수정 방법

### 1. 브랜치 생성
```bash
git checkout -b feat/update-vue-template
```

### 2. 템플릿 수정
```bash
cd templates/vuejs
# template.tmpl, config.json, metadata.json 수정
```

### 3. PR 생성
- PR 생성 시 CI가 자동으로 실행됩니다
- 변경된 템플릿 검증
- 마이그레이션 SQL 미리보기

### 4. 머지 및 자동 배포
- main 브랜치 머지 시 CD가 자동 실행
- 운영 MySQL에 직접 마이그레이션 적용
- 변경된 템플릿만 deprecated → active 처리

## 🚀 로컬 개발

### 마이그레이션 생성 테스트

**모든 템플릿:**
```bash
cd scripts
go run generate-migration.go --all --output /tmp/test.sql --version test
cat /tmp/test.sql
```

**특정 템플릿만:**
```bash
go run generate-migration.go \
  --changed '["../templates/vuejs","../templates/react"]' \
  --output /tmp/test.sql \
  --version test
```

**검증만 (SQL 생성 안 함):**
```bash
go run generate-migration.go --all --validate-only
```

## 📝 템플릿 파일 구조

### template.tmpl
Dockerfile 템플릿 파일 (Go template 문법 사용)

```dockerfile
FROM node:{{ .node_version }}-alpine

{{ if eq .package_manager "npm" }}
RUN npm ci
{{ else }}
RUN yarn install
{{ end }}

...
```

### metadata.json
템플릿 기본 정보

```json
{
  "name": "Vue.js",
  "description": "Vue.js 기반 SPA 애플리케이션",
  "categories": ["frontend"],
  "display_order": 1,
  "icon_name": "vue",
  "requires_git": true,
  "version": "1.0"
}
```

### config.json
템플릿 설정 (옵션, 포트, 환경변수)

```json
{
  "template_options": [
    {
      "name": "node_version",
      "label": "Node.js 버전",
      "type": "select",
      "options": ["18", "20", "22"],
      "default": "20"
    }
  ],
  "template_ports": [
    {
      "internal_port": 80,
      "network_type": "http",
      "description": "HTTP port"
    }
  ],
  "default_ports": [...],
  "default_env": [...]
}
```

## 🔄 CI/CD 워크플로우

### CI (Pull Request)
1. 변경된 템플릿 감지 (Git diff)
2. 템플릿 검증
3. 마이그레이션 SQL 생성 테스트
4. PR 코멘트에 미리보기 추가

### CD (Main Branch Push)
1. 변경된 템플릿 감지
2. 마이그레이션 SQL 생성
3. SSH로 운영 서버 접속
4. MySQL에 직접 마이그레이션 적용
5. 검증 (deprecated/active 상태 확인)

## 🛡️ 안전 장치

- **변경된 템플릿만 업데이트**: 불필요한 template_id 증가 방지
- **Deprecated 방식**: 기존 컨테이너는 이전 템플릿 계속 사용 가능
- **CI 검증 필수**: PR 승인 없이는 배포 불가
- **투명한 이력**: GitHub에서 모든 변경 추적 가능

## 📊 데이터베이스 스키마

### TEMPLATES 테이블
```sql
template_id     INT UNSIGNED AUTO_INCREMENT
name            VARCHAR(255)  -- 템플릿 이름
template_body   LONGTEXT      -- Dockerfile 템플릿
template_config JSON          -- 메타데이터 + 설정
status          ENUM('active', 'inactive', 'deprecated')
created_at      TIMESTAMP
updated_at      TIMESTAMP
```

### 배포 전략
1. **UPDATE**: 기존 템플릿 `status = 'deprecated'`
2. **INSERT**: 새 버전 템플릿 `status = 'active'`
3. 기존 컨테이너는 `deprecated` 템플릿 계속 사용
4. 신규 컨테이너는 `active` 템플릿만 선택 가능

## 🔗 관련 레포지토리

- [web-console-backend](https://github.com/swm-launchpad/web-console-backend) - LaunchPad Backend API
- [web-console-frontend](https://github.com/swm-launchpad/web-console-frontend) - LaunchPad Frontend
