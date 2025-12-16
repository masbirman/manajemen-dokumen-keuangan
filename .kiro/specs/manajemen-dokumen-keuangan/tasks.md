# Implementation Plan

## Phase 1: Project Setup & Infrastructure

- [x] 1. Setup Docker development environment

  - [x] 1.1 Create project directory structure (backend/, frontend/)

    - Create root docker-compose.yml
    - Create backend/Dockerfile with Go 1.22 and Air hot reload
    - Create frontend/Dockerfile with Node 20 and Vite
    - _Requirements: Tech Stack_

  - [x] 1.2 Initialize Goravel backend project

    - Install Goravel framework
    - Configure database connection for PostgreSQL
    - Setup .env file with Docker environment variables
    - _Requirements: Tech Stack_

  - [x] 1.3 Initialize Vue.js 3 frontend project

    - Create Vite + Vue 3 + TypeScript project
    - Install Tailwind CSS
    - Install Pinia for state management
    - Install Vue Router
    - Install Axios for API calls
    - _Requirements: Tech Stack_

  - [x] 1.4 Verify Docker environment works

    - Run docker-compose up and verify all services start
    - Verify frontend can reach backend API
    - Verify backend can connect to PostgreSQL
    - _Requirements: Tech Stack_

## Phase 2: Database & Models

- [x] 2. Create database migrations and models

  - [x] 2.1 Create database migrations

    - Create users table migration
    - Create unit_kerja table migration
    - Create pptk table migration
    - Create sumber_dana table migration
    - Create jenis_dokumen table migration
    - Create dokumen table migration
    - Create settings table migration
    - _Requirements: 1.1, 2.1, 3.1, 4.1, 5.1, 6.3, 9.1_

  - [x] 2.2 Create Go model structs

    - Create User model with relationships
    - Create UnitKerja model
    - Create PPTK model with UnitKerja relationship
    - Create SumberDana model
    - Create JenisDokumen model
    - Create Dokumen model with all relationships
    - Create Setting model
    - _Requirements: 1.1, 2.1, 3.1, 4.1, 5.1, 6.3, 9.1_

  - [x]\* 2.3 Write property test for master data CRUD

    - **Property 5: Master Data CRUD Consistency**
    - **Validates: Requirements 2.1, 3.1, 4.1, 4.2, 5.1, 5.2**

  - [x] 2.4 Create database seeder for initial Super Admin

    - Seed default Super Admin user
    - _Requirements: 1.1_

- [x] 3. Checkpoint - Verify database setup

  - Ensure all migrations run successfully
  - Ensure seeder creates Super Admin
  - Ask the user if questions arise

## Phase 3: Authentication & Authorization

- [x] 4. Implement authentication system

  - [x] 4.1 Create auth repository and service

    - Implement user lookup by username
    - Implement password hashing with bcrypt
    - Implement JWT token generation (access + refresh)
    - _Requirements: 10.1, 10.2_

  - [x] 4.2 Create auth controller and routes

    - POST /api/auth/login endpoint
    - POST /api/auth/logout endpoint
    - POST /api/auth/refresh endpoint
    - GET /api/auth/me endpoint
    - _Requirements: 10.1, 10.2_

  - [x] 4.3 Write property tests for authentication

    - **Property 21: Authentication Success**
    - **Property 22: Authentication Failure**
    - **Validates: Requirements 10.1, 10.2**

  - [x] 4.4 Create auth middleware

    - JWT validation middleware
    - Extract user from token and attach to context
    - _Requirements: 10.5_

  - [x] 4.5 Create role-based access middleware

    - Check user role against required role
    - Return 403 for insufficient permissions
    - _Requirements: 10.3, 10.4_

  - [x] 4.6 Write property test for RBAC

    - **Property 23: Role-Based Access Control**
    - **Validates: Requirements 10.3, 10.4**

  - [x] 4.7 Write property test for session expiration

    - **Property 24: Session Expiration Handling**
    - **Validates: Requirements 10.5**

- [x] 5. Checkpoint - Verify authentication works

  - Ensure login/logout works
  - Ensure JWT validation works
  - Ensure role middleware blocks unauthorized access
  - Ask the user if questions arise

## Phase 4: Master Data Management (Backend)

- [x] 6. Implement Unit Kerja CRUD

  - [x] 6.1 Create Unit Kerja repository

    - Implement Create, Read, Update, Delete operations

    - Implement GetAll with pagination
    - _Requirements: 2.1, 2.5_

  - [x] 6.2 Create Unit Kerja service

    - Implement business logic for CRUD
    - Implement validation
    - _Requirements: 2.1, 2.5_

  - [x] 6.3 Create Unit Kerja controller and routes

    - GET /api/unit-kerja
    - POST /api/unit-kerja
    - PUT /api/unit-kerja/:id
    - DELETE /api/unit-kerja/:id
    - _Requirements: 2.1, 2.5_

- [x] 7. Implement PPTK CRUD

  - [x] 7.1 Create PPTK repository

    - Implement CRUD operations
    - Implement GetByUnitKerja
    - _Requirements: 3.1, 3.2_

  - [x] 7.2 Create PPTK service

    - Implement business logic
    - Implement Unit Kerja assignment
    - _Requirements: 3.1, 3.2_

  - [x] 7.3 Create PPTK controller and routes

    - GET /api/pptk
    - POST /api/pptk
    - PUT /api/pptk/:id
    - DELETE /api/pptk/:id
    - GET /api/pptk/by-unit-kerja/:unitKerjaId
    - _Requirements: 3.1, 3.2_

- [x] 8. Implement Sumber Dana and Jenis Dokumen CRUD

  - [x] 8.1 Create Sumber Dana repository, service, controller

    - Full CRUD implementation

    - Referential integrity check on delete
    - _Requirements: 4.1, 4.2, 4.3_

  - [x] 8.2 Create Jenis Dokumen repository, service, controller

    - Full CRUD implementation
    - Referential integrity check on delete
    - _Requirements: 5.1, 5.2, 5.3_

  - [x] 8.3 Write property test for referential integrity

    - **Property 6: Referential Integrity on Delete**
    - **Validates: Requirements 4.3, 5.3**

- [x] 9. Checkpoint - Verify master data CRUD

  - Ensure all CRUD endpoints work
  - Ensure referential integrity is enforced
  - Ask the user if questions arise

## Phase 5: Excel Import/Export (Backend)

- [x] 10. Implement Excel functionality

  - [x] 10.1 Create Excel service using excelize

    - Implement template generation for Unit Kerja
    - Implement template generation for PPTK
    - Implement Excel parsing and validation
    - _Requirements: 2.2, 2.3, 2.4, 3.3, 3.4, 3.5_

  - [x] 10.2 Add Excel endpoints to Unit Kerja controller

    - GET /api/unit-kerja/template
    - POST /api/unit-kerja/import
    - _Requirements: 2.2, 2.3, 2.4_

  - [x] 10.3 Add Excel endpoints to PPTK controller

    - GET /api/pptk/template
    - POST /api/pptk/import
    - _Requirements: 3.3, 3.4, 3.5_

  - [x] 10.4 Write property tests for Excel import

    - **Property 7: Excel Import Creates All Valid Records**
    - **Property 8: Excel Import Error Reporting**
    - **Validates: Requirements 2.3, 2.4, 3.4, 3.5**

## Phase 6: User Management (Backend)

- [x] 11. Implement User management

  - [x] 11.1 Create User repository and service

    - Implement CRUD operations
    - Implement Operator assignment to PPTK and Unit Kerja
    - Implement user deactivation
    - _Requirements: 1.1, 1.2, 1.4, 1.5_

  - [x] 11.2 Create User controller and routes

    - GET /api/users
    - POST /api/users
    - GET /api/users/:id
    - PUT /api/users/:id
    - DELETE /api/users/:id (deactivate)
    - _Requirements: 1.1, 1.2, 1.4, 1.5_

  - [x] 11.3 Write property tests for user management

    - **Property 1: User Creation Persistence**
    - **Property 2: Operator Assignment Integrity**
    - **Property 4: User Deactivation Blocks Login**
    - **Validates: Requirements 1.1, 1.2, 1.5**

## Phase 7: File Upload (Backend)

- [x] 12. Implement file upload functionality

  - [x] 12.1 Create file service

    - Implement file upload to storage
    - Implement file type validation (PDF only for documents)
    - Implement image validation for avatars
    - _Requirements: 1.3, 3.6, 6.4_

  - [x] 12.2 Add avatar upload endpoints

    - POST /api/users/:id/avatar
    - POST /api/pptk/:id/avatar
    - _Requirements: 1.3, 3.6_

  - [x] 12.3 Write property tests for file upload

    - **Property 3: Avatar Upload Persistence**
    - **Property 11: PDF File Type Validation**
    - **Validates: Requirements 1.3, 3.6, 6.4**

- [x] 13. Checkpoint - Verify backend features

  - Ensure Excel import/export works
  - Ensure user management works
  - Ensure file upload works
  - Ask the user if questions arise

## Phase 8: Document Management (Backend)

- [x] 14. Implement Document CRUD

  - [x] 14.1 Create Dokumen repository

    - Implement Create operation
    - Implement GetAll with filters (unit_kerja, date_range)
    - Implement role-based filtering (Operator sees own, Admin sees all)
    - _Requirements: 6.3, 8.1, 8.2, 8.3, 8.4_

  - [x] 14.2 Create Dokumen service

    - Implement document creation with file upload
    - Implement validation for required fields
    - Implement filter logic
    - _Requirements: 6.3, 6.5, 8.1, 8.2, 8.3, 8.4_

  - [x] 14.3 Create Dokumen controller and routes

    - GET /api/dokumen (with filters)
    - POST /api/dokumen
    - GET /api/dokumen/:id
    - GET /api/dokumen/:id/file
    - _Requirements: 6.3, 8.1, 8.5_

  - [x] 14.4 Write property tests for document management

    - **Property 10: Document Creation Completeness**
    - **Property 12: Document Validation Errors**
    - **Property 16: Admin Document Visibility**
    - **Property 17: Operator Document Isolation**
    - **Property 18: Unit Kerja Filter Accuracy**
    - **Property 19: Date Range Filter Accuracy**
    - **Validates: Requirements 6.3, 6.5, 8.1, 8.2, 8.3, 8.4**

## Phase 9: Settings (Backend)

- [x] 15. Implement Settings management

  - [x] 15.1 Create Settings repository, service, controller

    - GET /api/settings
    - PUT /api/settings
    - _Requirements: 9.1, 9.2_

  - [x] 15.2 Write property test for settings

    - **Property 20: Settings Persistence**
    - **Validates: Requirements 9.2**

- [x] 16. Checkpoint - Verify all backend APIs

  - Ensure all endpoints work correctly

  - Ensure all property tests pass
  - Ask the user if questions arise

## Phase 10: Frontend - Core Setup

- [x] 17. Setup frontend core infrastructure
  - [x] 17.1 Create API service with Axios
    - Setup base URL from environment
    - Setup request/response interceptors
    - Handle JWT token in headers
    - Handle 401 redirect to login
    - _Requirements: 10.1, 10.5_
  - [x] 17.2 Create Pinia auth store
    - Store user data and tokens
    - Implement login/logout actions
    - Implement token refresh
    - _Requirements: 10.1_
  - [x] 17.3 Create Vue Router with guards
    - Setup route definitions
    - Implement auth guard
    - Implement role-based route guards
    - _Requirements: 10.3, 10.4_
  - [x] 17.4 Create main layout components
    - Create MainLayout.vue with sidebar and navbar
    - Create Sidebar.vue with role-based menu
    - Create Navbar.vue with user info and logout
    - _Requirements: 10.3, 10.4_

## Phase 11: Frontend - Reusable Components

- [x] 18. Create reusable UI components

  - [x] 18.1 Create form components

    - InputField.vue with validation display
    - Dropdown.vue with search and default value
    - CurrencyInput.vue with formatting
    - FileUpload.vue with drag-drop and preview
    - _Requirements: 6.1, 6.2, 6.6_

  - [x] 18.2 Write property test for currency input

    - **Property 13: Currency Value Handling**

    - **Validates: Requirements 6.6**

  - [x] 18.3 Create data display components

    - DataTable.vue with sorting, filtering, pagination
    - Modal.vue for dialogs
    - Toast notifications
    - _Requirements: 8.1_

## Phase 12: Frontend - Authentication Pages

- [x] 19. Create authentication pages

  - [x] 19.1 Create LoginView.vue

    - Login form with validation
    - Error message display
    - Redirect after successful login
    - _Requirements: 10.1, 10.2_

## Phase 13: Frontend - Master Data Pages

- [x] 20. Create master data management pages

  - [x] 20.1 Create UnitKerjaView.vue

    - Data table with CRUD operations
    - Excel template download button
    - Excel import with file upload
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_

  - [x] 20.2 Create PPTKView.vue

    - Data table with CRUD operations
    - Unit Kerja dropdown for assignment
    - Avatar upload
    - Excel template download and import
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6_

  - [x] 20.3 Create SumberDanaView.vue

    - Data table with CRUD operations
    - Delete confirmation with referential integrity warning

    - _Requirements: 4.1, 4.2, 4.3_

  - [x] 20.4 Create JenisDokumenView.vue

    - Data table with CRUD operations
    - Delete confirmation with referential integrity warning
    - _Requirements: 5.1, 5.2, 5.3_

## Phase 14: Frontend - User Management

- [x] 21. Create user management page

  - [x] 21.1 Create ManajemenUserView.vue

    - Data table with user list
    - Create/edit user form with role selection
    - Unit Kerja and PPTK assignment for Operators
    - Avatar upload
    - Activate/deactivate user
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_

## Phase 15: Frontend - Document Input

- [x] 22. Create document input functionality

  - [x] 22.1 Create InputDokumenView.vue

    - Form with all required fields

    - Dropdowns with default values based on Operator assignment
    - PDF file upload
    - Form validation
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6_

  - [x] 22.2 Write property test for operator defaults

    - **Property 9: Operator Default Selection**
    - **Validates: Requirements 6.1, 6.2**

  - [x] 22.3 Create DocumentScanner.vue component

    - Camera access using MediaDevices API
    - Capture image functionality
    - Multi-page capture support
    - Convert images to PDF using jsPDF
    - Responsive: show only on mobile browsers
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5_

  - [x] 22.4 Write property tests for scanner

    - **Property 14: Image to PDF Conversion**
    - **Property 15: Multi-page PDF Combination**
    - **Validates: Requirements 7.2, 7.3**

## Phase 16: Frontend - Document List

- [x] 23. Create document list page

  - [x] 23.1 Create ListDokumenView.vue

    - Data table with all document fields
    - Filter by Unit Kerja dropdown
    - Filter by date range
    - Role-based data display (Operator sees own, Admin sees all)
    - Click row to view detail and PDF preview
    - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5_

## Phase 17: Frontend - Settings & Dashboard

- [x] 24. Create settings and dashboard pages

  - [x] 24.1 Create PengaturanView.vue

    - Display configurable settings
    - Edit settings form
    - _Requirements: 9.1, 9.2, 9.3_

  - [x] 24.2 Create DashboardView.vue

    - Summary statistics
    - Recent documents
    - Quick actions based on role
    - _Requirements: General UX_

- [x] 25. Final Checkpoint - Full system verification

  - Ensure all features work end-to-end
  - Ensure all tests pass
  - Ask the user if questions arise
