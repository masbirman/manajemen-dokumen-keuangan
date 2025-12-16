# Requirements Document

## Introduction

Sistem Manajemen Dokumen Keuangan adalah aplikasi full web (responsive) yang memungkinkan pengelolaan dokumen keuangan secara digital. Sistem ini mendukung tiga peran pengguna (Super Admin, Admin, Operator) dengan fitur input dokumen, manajemen data master, dan scanner dokumen via browser camera untuk perangkat mobile.

## Tech Stack

- **Backend**: Goravel (Go framework dengan struktur Laravel-like)
- **Frontend**: Vue.js 3 + Vite + Tailwind CSS
- **Database**: PostgreSQL
- **File Storage**: Local/Cloud storage via Goravel filesystem
- **Excel Processing**: excelize (Go library)
- **PDF Processing**: jsPDF / pdf-lib (JavaScript)
- **Camera/Scanner**: Browser MediaDevices API + vue-web-cam

## Glossary

- **Unit_Kerja**: Entitas organisasi yang menjadi sumber penugasan PPTK dan Operator
- **PPTK**: Pejabat Pelaksana Teknis Kegiatan yang bertanggung jawab atas dokumen keuangan
- **Sumber_Dana**: Kategori asal dana yang digunakan dalam dokumen keuangan
- **Jenis_Dokumen**: Klasifikasi tipe dokumen keuangan
- **Operator**: Pengguna yang melakukan input dokumen, di-assign ke PPTK dan Unit Kerja tertentu
- **Admin**: Pengguna yang mengelola data master dan melihat semua hasil inputan
- **Super_Admin**: Pengguna dengan akses penuh termasuk manajemen user dan pengaturan sistem
- **Avatar**: Gambar profil pengguna atau PPTK
- **Scanner_Dokumen**: Fitur untuk memindai dokumen fisik menjadi file PDF menggunakan kamera perangkat (khusus untuk browser mobile/responsive web)

## Requirements

### Requirement 1: Manajemen User

**User Story:** As a Super Admin, I want to manage user accounts, so that I can control access to the system.

#### Acceptance Criteria

1. WHEN a Super Admin creates a new user THEN the System SHALL store the user with username, password, role, and optional avatar
2. WHEN a Super Admin assigns an Operator to a PPTK and Unit Kerja THEN the System SHALL link the Operator to both entities
3. WHEN a user uploads an avatar THEN the System SHALL store and display the avatar image on the user profile
4. WHEN a Super Admin updates user information THEN the System SHALL persist the changes immediately
5. WHEN a Super Admin deactivates a user THEN the System SHALL prevent that user from logging in

### Requirement 2: Manajemen Data Unit Kerja

**User Story:** As an Admin, I want to manage Unit Kerja data with Excel import capability, so that I can efficiently maintain organizational structure.

#### Acceptance Criteria

1. WHEN an Admin creates a Unit Kerja THEN the System SHALL store the Unit Kerja with unique identifier and name
2. WHEN an Admin requests an Excel template for Unit Kerja THEN the System SHALL provide a downloadable template file with correct column headers
3. WHEN an Admin imports Unit Kerja data from Excel THEN the System SHALL validate and create multiple Unit Kerja records from the file
4. WHEN an Admin imports Excel with invalid data THEN the System SHALL report specific validation errors with row numbers
5. WHEN an Admin updates a Unit Kerja THEN the System SHALL persist the changes and maintain referential integrity with assigned PPTK

### Requirement 3: Manajemen Data PPTK

**User Story:** As an Admin, I want to manage PPTK data with Excel import capability, so that I can efficiently maintain PPTK records.

#### Acceptance Criteria

1. WHEN an Admin creates a PPTK THEN the System SHALL store the PPTK with name, Unit Kerja assignment, and optional avatar
2. WHEN an Admin assigns a PPTK to a Unit Kerja THEN the System SHALL link the PPTK to that Unit Kerja
3. WHEN an Admin requests an Excel template for PPTK THEN the System SHALL provide a downloadable template file with correct column headers
4. WHEN an Admin imports PPTK data from Excel THEN the System SHALL validate and create multiple PPTK records from the file
5. WHEN an Admin imports Excel with invalid PPTK data THEN the System SHALL report specific validation errors with row numbers
6. WHEN a PPTK uploads an avatar THEN the System SHALL store and display the avatar image on the PPTK profile

### Requirement 4: Manajemen Data Sumber Dana

**User Story:** As an Admin, I want to manage Sumber Dana data, so that I can categorize funding sources for documents.

#### Acceptance Criteria

1. WHEN an Admin creates a Sumber Dana THEN the System SHALL store the Sumber Dana with unique identifier and name
2. WHEN an Admin updates a Sumber Dana THEN the System SHALL persist the changes
3. WHEN an Admin deletes a Sumber Dana that is referenced by documents THEN the System SHALL prevent deletion and display an error message

### Requirement 5: Manajemen Data Jenis Dokumen

**User Story:** As an Admin, I want to manage Jenis Dokumen data, so that I can categorize document types.

#### Acceptance Criteria

1. WHEN an Admin creates a Jenis Dokumen THEN the System SHALL store the Jenis Dokumen with unique identifier and name
2. WHEN an Admin updates a Jenis Dokumen THEN the System SHALL persist the changes
3. WHEN an Admin deletes a Jenis Dokumen that is referenced by documents THEN the System SHALL prevent deletion and display an error message

### Requirement 6: Input Dokumen Keuangan

**User Story:** As an Operator, I want to input financial documents with pre-filled defaults based on my assignment, so that I can efficiently record documents.

#### Acceptance Criteria

1. WHEN an Operator opens the input form THEN the System SHALL pre-select the Unit Kerja dropdown based on the Operator's assignment
2. WHEN an Operator opens the input form THEN the System SHALL pre-select the PPTK dropdown based on the Operator's assignment
3. WHEN an Operator submits a document THEN the System SHALL store Unit Kerja, PPTK, Jenis Dokumen, Sumber Dana, Nilai, Uraian, and uploaded PDF file
4. WHEN an Operator uploads a PDF file THEN the System SHALL validate the file type and store the file securely
5. WHEN an Operator submits a document with missing required fields THEN the System SHALL display validation errors for each missing field
6. WHEN an Operator inputs Nilai THEN the System SHALL accept numeric values and format them as currency

### Requirement 7: Scanner Dokumen (Responsive Web)

**User Story:** As an Operator using mobile browser, I want to scan physical documents using device camera, so that I can digitize paper documents directly from the web app.

#### Acceptance Criteria

1. WHEN an Operator accesses the scanner feature on mobile browser THEN the System SHALL request camera permission and activate document capture interface
2. WHEN an Operator captures a document image via browser THEN the System SHALL process the image and convert it to PDF format
3. WHEN an Operator scans multiple pages via browser THEN the System SHALL combine all pages into a single PDF file
4. WHEN an Operator completes scanning THEN the System SHALL attach the generated PDF to the document input form
5. WHEN an Operator accesses the web app on desktop browser THEN the System SHALL show only file upload option without scanner feature

### Requirement 8: Tabel Hasil Inputan

**User Story:** As an Admin or Super Admin, I want to view all inputted documents in a table, so that I can monitor and review financial records.

#### Acceptance Criteria

1. WHEN an Admin views the hasil inputan table THEN the System SHALL display all documents with Unit Kerja, PPTK, Jenis Dokumen, Sumber Dana, Nilai, Uraian, and submission date
2. WHEN an Operator views the hasil inputan table THEN the System SHALL display only documents submitted by that Operator
3. WHEN a user filters the table by Unit Kerja THEN the System SHALL display only documents matching the selected Unit Kerja
4. WHEN a user filters the table by date range THEN the System SHALL display only documents within the specified dates
5. WHEN a user clicks on a document row THEN the System SHALL display the document details and allow PDF preview

### Requirement 9: Pengaturan Sistem

**User Story:** As a Super Admin, I want to configure system settings, so that I can customize the application behavior.

#### Acceptance Criteria

1. WHEN a Super Admin accesses pengaturan THEN the System SHALL display configurable system parameters
2. WHEN a Super Admin updates a setting THEN the System SHALL persist the change and apply it immediately
3. WHEN a Super Admin updates application branding THEN the System SHALL reflect the changes across the application

### Requirement 10: Autentikasi dan Otorisasi

**User Story:** As a user, I want to securely login and access features based on my role, so that the system remains secure.

#### Acceptance Criteria

1. WHEN a user provides valid credentials THEN the System SHALL authenticate and create a session
2. WHEN a user provides invalid credentials THEN the System SHALL reject login and display an error message
3. WHEN an Operator attempts to access Admin features THEN the System SHALL deny access and display unauthorized message
4. WHEN an Admin attempts to access Super Admin features THEN the System SHALL deny access and display unauthorized message
5. WHEN a user session expires THEN the System SHALL redirect to login page
