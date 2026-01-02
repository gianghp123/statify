# Statify Frontend Design & Prompting Guide

This document provides a comprehensive guide for designing and implementing the frontend for **Statify**, a static site hosting management platform. It includes UI/UX requirements, API integration details, and sample prompts for AI-assisted development.

## 1. Project Overview
Statify allows users to host static websites by uploading ZIP files. The system handles subdomains, deployments, and project management.

## 2. Design System (Aesthetics)
- **Theme**: Modern, clean, and developer-centric.
- **Vibe**: Professional, reliable, and high-tech (similar to Vercel, Netlify, or Supabase).
- **Core Colors**:
    - Primary: Deep Indigo or Sleek Blue.
    - Status Colors: Green (Ready), Yellow (Queued/Processing), Red (Failed).
- **Typography**: Modern sans-serif (e.g., Inter, Montserrat).

## 3. UI Screens & Features

### A. Authentication
- **Login/Register**: Standard email/password forms.
- **GitHub OAuth**: "Continue with GitHub" button.
- **User Profile**: "Me" endpoint to display current user info.

### B. Project Dashboard (List)
- **Feature**: Display a list of projects with pagination.
- **Data per Card**: Project Name, Subdomain (link), Status of current deployment, Created date.
- **Action**: "Create New Project" button.

### C. Create Project
- **Form**: Name, Subdomain (availability check encouraged).
- **Upload**: Dropzone for ZIP file.

### D. Project Details & Deployment History
- **Overview**: Show project stats (created at, active deployment).
- **History Table**: List deployments with columns: ID, Status, Created At, Finished At.
- **Logs/Errors**: If a deployment failed, show the `validation_error`.

## 4. API Integration Details

### Base URL: `/api/v1`

### Data Structures
- **Paginated Response**:
    ```json
    {
      "code": 200,
      "message": "...",
      "data": [...],
      "pagination": {
        "total_count": 100,
        "page": 1,
        "limit": 10
      }
    }
    ```
- **Error Response**:
    ```json
    {
      "code": 404,
      "message": "Not Found"
    }
    ```

### Endpoints
| Module | Method | Path | Description |
| :--- | :--- | :--- | :--- |
| **Auth** | POST | `/auth/register` | User registration |
| | POST | `/auth/login` | User login (returns JWT) |
| | GET | `/auth/me` | Get current user info |
| **Project**| GET | `/projects` | List projects (Paginated) |
| | POST | `/projects` | Create project |
| | GET | `/projects/:project_id` | Get project details |
| **Deployment**| GET | `/projects/:project_id/deployments` | List history (Paginated) |
| | POST | `/projects/:project_id/deployments` | Create new deployment |
| | GET | `/projects/:project_id/deployments/:id` | Get deployment status |

## 5. Sample Prompt for Frontend Generation

> [!TIP]
> Use the following prompt with an AI coding assistant (like Antigravity) to generate the UI components.

**Prompt:**
"Build a premium React (Vite) dashboard for 'Statify', a static site hosting platform. Use Vanilla CSS with a glassmorphism aesthetic and a dark mode theme. The dashboard should include:
1. A sidebar for navigation (Dashboard, Projects, Settings).
2. A main project list page using cards. Each card displays the project name, a clickable subdomain link, and a status badge (Ready, Failed, Queued).
3. A 'Create Project' modal with a file upload dropzone for ZIP files.
4. An 'Activity' section showing recent deployment history in a table.
Ensure the layout is fully responsive and uses the 'Inter' font. Use sleek animations for transitions."

## 6. Development Checklist
- [ ] Setup JWT authentication interceptor.
- [ ] Implement global loading states for deployments.
- [ ] Handle 401 Unauthorized errors by redirecting to login.
- [ ] Implement polling or WebSockets for deployment status updates.
