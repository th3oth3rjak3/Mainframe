# Mainframe

A personal mono-repository containing multiple applications for my daily workflow and productivity.

## Overview

This project is a collection of applications I'm building for my own use. It's a work in progress with new features and apps being added as I identify needs in my personal workflow.

## Applications

*Note: This is an evolving project. Applications listed may be in various stages of development.*

- **Auth System** - Centralized authentication with role-based access control
- **Recipes** - Recipe management and organization *(in development)*
- **Lists** - Personal list management for shopping, gifts, and tasks *(planned)*
- **Calendar** - Personal and shared calendar application *(planned)*
- **Password Manager** - Secure password storage and management *(planned)*

## Tech Stack

- **Backend**: ASP.NET Core, Entity Framework Core
- **Database**: PostgreSQL / SQLite
- **Authentication**: Cookie-based with database sessions
- **Password Hashing**: Argon2id

## Getting Started

```bash
# Clone the repository
git clone https://github.com/yourusername/mainframe.git
cd mainframe

# Restore dependencies
dotnet restore

# Run database migrations
dotnet ef database update

# Start the application
dotnet run
```

## Development Status

🚧 **This project is in active development for personal use.** 

Features are implemented as needed, and the scope may change based on my requirements. Feel free to explore the code, but expect frequent changes and incomplete features.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
