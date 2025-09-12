# Frontend Template Standards

This document describes the standardized configurations for all frontend templates.

## Package.json Standards

### Required Scripts

- `dev`: `next dev`
- `build`: `next build`
- `lint`: `next lint`
- `lint:fix`: `next lint --fix`
- `type-check`: `tsc --noEmit`
- `test:coverage`: `jest --coverage`
- `format`: `prettier --write .`
- `clean`: `rm -rf .next out dist`
- `start`: `next start`
- `test`: `jest`
- `test:watch`: `jest --watch`
- `format:check`: `prettier --check .`

### Required Dependencies

- `next`: `{{.Versions.NextJS}}`
- `react`: `{{.Versions.React}}`
- `react-dom`: `{{.Versions.React}}`
- `@types/node`: `^20.10.0`
- `@types/react-dom`: `^18.2.0`
- `postcss`: `^8.4.0`
- `class-variance-authority`: `^0.7.0`
- `@radix-ui/react-slot`: `^1.0.0`
- `typescript`: `^5.3.0`
- `tailwindcss`: `^3.4.0`
- `clsx`: `^2.0.0`
- `lucide-react`: `^0.300.0`
- `tailwindcss-animate`: `^1.0.7`
- `tailwind-merge`: `^2.2.0`
- `@types/react`: `^18.2.0`
- `autoprefixer`: `^10.4.0`

### Required Dev Dependencies

- `prettier`: `^3.1.0`
- `prettier-plugin-tailwindcss`: `^0.5.0`
- `@types/jest`: `^29.5.0`
- `eslint`: `^8.55.0`
- `@typescript-eslint/eslint-plugin`: `^6.15.0`
- `jest`: `^29.7.0`
- `jest-environment-jsdom`: `^29.7.0`
- `@testing-library/react`: `^14.1.0`
- `@testing-library/jest-dom`: `^6.1.0`
- `@testing-library/user-event`: `^14.5.0`
- `eslint-config-next`: `{{.Versions.NextJS}}`
- `@typescript-eslint/parser`: `^6.15.0`

### Engine Requirements

- `node`: `>=22.0.0`
- `npm`: `>=10.0.0`

## Template-Specific Configurations

### nextjs-app

- **Port**: 3000
- **Additional Dependencies**:
  - `@radix-ui/react-dialog`: `^1.0.0`
  - `@radix-ui/react-dropdown-menu`: `^2.0.0`
  - `@radix-ui/react-toast`: `^1.1.0`

### nextjs-home

- **Port**: 3001
- **Additional Dependencies**:
  - `@radix-ui/react-accordion`: `^1.1.0`
  - `@radix-ui/react-navigation-menu`: `^1.1.0`
  - `framer-motion`: `^10.16.0`
  - `react-intersection-observer`: `^9.5.0`

### nextjs-admin

- **Port**: 3002
- **Additional Dependencies**:
  - `@radix-ui/react-checkbox`: `^1.0.0`
  - `@radix-ui/react-toast`: `^1.1.0`
  - `zod`: `^3.22.0`
  - `date-fns`: `^3.0.0`
  - `@radix-ui/react-dropdown-menu`: `^2.0.0`
  - `@radix-ui/react-switch`: `^1.0.0`
  - `@radix-ui/react-tabs`: `^1.0.0`
  - `@radix-ui/react-tooltip`: `^1.0.0`
  - `@tanstack/react-table`: `^8.11.0`
  - `react-hook-form`: `^7.48.0`
  - `@hookform/resolvers`: `^3.3.0`
  - `recharts`: `^2.8.0`
  - `@radix-ui/react-select`: `^2.0.0`

## Configuration Files

All frontend templates must include the following standardized configuration files:

- `package.json.tmpl`
- `tsconfig.json.tmpl`
- `.eslintrc.json.tmpl`
- `.prettierrc.tmpl`
- `vercel.json.tmpl`
- `tailwind.config.js.tmpl`
- `next.config.js.tmpl`
- `postcss.config.js.tmpl`
- `jest.config.js.tmpl`
- `jest.setup.js.tmpl`

## Vercel Deployment Standards

All frontend templates are configured for Vercel deployment with:

- **Framework**: Next.js
- **Build Command**: `npm run build`
- **Dev Command**: `npm run dev`
- **Install Command**: `npm install`
- **Security Headers**: Configured for production security
- **Environment Variables**: Standardized naming conventions

## Validation Rules

Templates are validated against the following rules:

1. **Package.json Consistency**: All required scripts, dependencies, and engines must be present
2. **TypeScript Configuration**: Standardized compiler options and path mappings
3. **ESLint Configuration**: Consistent linting rules across all templates
4. **Prettier Configuration**: Uniform code formatting settings
5. **Vercel Compatibility**: Proper deployment configuration
6. **Security Standards**: Required security headers and configurations
