# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ModelGate Web is the admin frontend for the ModelGate LLM API gateway platform. Built on Vue 3 + TypeScript + NaiveUI with a monorepo structure containing shared packages.

## Ignored Directories
Do not read, analyze, or modify files in the following directories:
- node_modules/
- dist/

## Common Commands

```bash
# Install dependencies
pnpm install

# Development server (port 9527)
pnpm dev

# Development with production env
pnpm dev:prod

# Build for production
pnpm build

# Build for test environment
pnpm build:test

# Lint
pnpm lint

# Type check
pnpm typecheck

# Preview production build
pnpm preview

# Generate routes (elegant-router)
pnpm gen-route

# Git commit with conventional commits
pnpm commit        # English
pnpm commit:zh     # Chinese
```

## Environment

- Node.js: >= 18.20.0
- pnpm: >= 8.7.0 (required, uses workspace features)

Key environment variables (`.env`):
- `VITE_PUBLIC_GRPC_SERVICE_URL` - gRPC backend URL (default: `http://localhost:8080/api`)
- `VITE_AUTH_ROUTE_MODE` - Route mode: `static` or `dynamic`
- `VITE_ROUTER_HISTORY_MODE` - Vue Router mode: `hash` or `history`
- `VITE_STORAGE_PREFIX` - LocalStorage prefix (default: `MG_`)

## Architecture

### Directory Structure

```
src/
├── grpc.ts              # gRPC client setup with interceptors
├── main.ts              # App entry point
├── App.vue              # Root component
├── router/              # Routing (elegant-router generated)
│   ├── elegant/         # Generated route definitions
│   ├── guard/           # Route guards (auth, permissions)
│   └── routes/          # Static route configs
├── store/               # Pinia state management
│   └── modules/         # auth, app, route, tab, theme
├── views/               # Page components
│   ├── relay/           # Relay management (model, provider, pricing)
│   ├── manage/          # System management (user, role, menu, permission)
│   ├── usage/           # Usage tracking (ledger, request)
│   ├── user/            # User settings (account, api-key)
│   └── _builtin/        # Error pages, login
├── layouts/             # Layout components (base, blank)
├── components/          # Shared components (common, custom, advanced)
├── hooks/               # Composables
│   ├── common/          # Generic hooks (table, loading, request)
│   └── business/        # Domain-specific hooks (auth)
├── typings/             # TypeScript types
│   └── proto/           # Generated protobuf types
├── locales/             # i18n translations
├── utils/               # Utility functions
├── constants/           # Constants and enums
├── theme/               # Theme configuration
└── styles/              # SCSS/CSS styles
└── plugins/             # Vue plugins

packages/                # Monorepo shared packages
├── axios/               # Axios wrapper
├── hooks/               # Shared hooks (useTable, useBoolean, etc.)
├── utils/               # Shared utilities (storage, crypto, nanoid)
├── materials/           # UI components (admin-layout, page-tab)
├── color/               # Color utilities
├── uno-preset/          # UnoCSS preset
├── scripts/             # Build scripts (gen-route, git-commit)
└── ofetch/              # Fetch wrapper
└── alova/               # Alova request library
```

### gRPC Communication

Uses ConnectRPC for gRPC-Web. Clients defined in `src/grpc.ts`:

```typescript
// Three service clients
authServiceClient    // Auth: login, register, token refresh
systemServiceClient  // System: users, roles, permissions, menus
relayServiceClient   // Relay: models, providers, pricing, accounts
```

Auth interceptor automatically:
- Adds `Authorization: Bearer <token>` header
- Handles token refresh on `Unauthenticated` errors
- Redirects to login on refresh failure

Protobuf types are generated in `src/typings/proto/` from `/proto` definitions using `buf generate`.

### Key Patterns

#### useTable Hook (packages/hooks/src/use-table.ts)

Standard pattern for list pages:

```typescript
const {
  columns,        // Table column definitions
  columnChecks,   // Column visibility controls
  data,           // Table data
  getData,        // Fetch data
  getDataByPage,  // Fetch with pagination reset
  loading,        // Loading state
  searchParams,   // Search parameters
  resetSearchParams,
  handleSorterChange
} = useTable({
  apiFn: relayServiceClient.getModelList,  // gRPC method
  showTotal: true,
  apiParams: { current: 1, size: 15, ... },  // Initial params
  columns: () => [...]  // Column factory
});
```

#### useTableOperate Hook

For CRUD operations:

```typescript
const {
  drawerVisible,   // Modal/drawer visibility
  operateType,     // 'add' | 'edit'
  editingData,     // Current editing row
  handleAdd,       // Open add modal
  handleEdit,      // Open edit modal with row data
  checkedRowKeys,  // Selected row IDs
  onBatchDeleted,  // Callback after batch delete
  onDeleted        // Callback after single delete
} = useTableOperate<Model>(data, getData);
```

#### View Structure Convention

Each feature module follows this pattern:

```
views/relay/model/
├── index.vue              # Main page (table + search)
└── modules/
    ├── model-search.vue   # Search form component
    └── model-operate-modal.vue  # Add/Edit drawer/modal
```

### Routing

Uses `@elegant-router/vue` - routes auto-generated from file structure:

- Route names derived from directory: `views/relay/model/` → `relay_model`
- Layout specified in route component: `layout.base$view.relay_model`
- Dynamic routes loaded from backend via `store/modules/route`

Route guards in `router/guard/`:
- Auth guard: checks token, redirects to login
- Permission guard: checks route permissions
- Title guard: sets document title

### State Management (Pinia)

Key store modules:
- `auth`: Token, user info, login/logout
- `app`: Theme, sidebar, locale settings
- `route`: Dynamic route permissions
- `tab`: Page tabs state
- `theme`: Theme preferences

### Styling

- UnoCSS for utility classes (Tailwind-like)
- SCSS for component styles with global mixins (`@styles/scss/global.scss`)
- NaiveUI component library with theme customization

### i18n

Translations in `src/locales/langs/`. Use `$t()` function:

```typescript
import { $t } from '@/locales';
$t('page.relay.model.title')  // Translation key
```

## Development Workflow

1. Adding a new page:
   - Create directory in `views/`
   - Create `index.vue` with table/search pattern
   - Create `modules/` for search/modal components
   - Run `pnpm gen-route` to regenerate routes
   - Add translations in `locales/langs/`

2. Adding gRPC methods:
   - Update proto files in `/proto`
   - Run `buf generate` (from root directory)
   - Use new client methods in views

3. Working with packages:
   - Packages are workspace dependencies (`workspace:*`)
   - Changes require `pnpm install` to update links