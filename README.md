🎬 MagicStreamMovies
A modern, full-stack movie streaming platform built with React and Go

🔗Webiste Hosted link: **https://magic-stream-movies-alpha.vercel.app/**


🎥 Features:
🎬 Movie Catalog - Browse and discover movies with detailed information
🔍 Search & Filter - Find your favorite movies quickly
📱 Responsive Design - Seamless experience across all devices
⚡ Real-time Streaming - Fast and reliable video streaming



🚀 Tech Stack
Frontend:
React.js - UI framework
JavaScript - Core programming language
Vercel - Deployment platform

Backend:
Go (Golang) - High-performance server
RESTful API - Clean API architecture

📁 Project Structure
MagicStreamMovies/
├── Client/                 # Frontend React application
├── Server/
│   └── MagicStreamMoviesServer/   # Backend Go server
└── magic-seed-data/       # Database seed data

🛠️ Installation
Prerequisites:
Node.js (v14 or higher)
Go (v1.19 or higher)
npm or yarn

Clone the Repository:
bashgit clone https://github.com/Kudwa-Abhishek/MagicStreamMovies.git
cd MagicStreamMovies

Frontend Setup:
bashcd Client
npm install
npm start
The frontend will run on http://localhost:3000

Backend Setup:
bashcd Server/MagicStreamMoviesServer
go mod download
go run main.go
The backend server will start on http://localhost:8080

🔧 Configuration
Environment Variables

Frontend (.env)
envREACT_APP_API_URL=http://localhost:8080

Backend (.env)
envPORT=8080
DATABASE_URL=your_database_url
JWT_SECRET=your_jwt_secret

📦 Database Setup
Seed the database with initial data:
bash# Navigate to the seed data directory
cd magic-seed-data
# Follow the instructions in the seed data folder

🌐 Deployment
Frontend (Vercel)
The frontend is automatically deployed to Vercel on every push to the main branch.
Visit: https://magic-stream-movies-alpha.vercel.app

Backend
Deploy the Go server to your preferred hosting platform:
Render
Railway
AWS EC2

🤝 Contributing
Contributions are welcome!

📝 API Documentation
Endpoints
Movies

GET /api/movies - Get all movies
GET /api/movies/:id - Get movie by ID
POST /api/movies - Create new movie (admin only)
PUT /api/movies/:id - Update movie (admin only)
DELETE /api/movies/:id - Delete movie (admin only)

Authentication
POST /api/auth/register - Register new user
POST /api/auth/login - User login
GET /api/auth/profile - Get user profile

🔐 Security
🔑 JWT-based authentication
🔒 Password hashing with bcrypt
🌐 CORS configuration
✅ Input validation and sanitization
🛡️ SQL injection prevention

SIMPLE AI FEATURE:
Used OPEN-AI with prompt such that a given review in natural language can be converted into one of these single-word reviews with the help of a straight prompt:
a) Excellent
b) Good
c) OK
d) Terrible
e) No review (Blank)

Used OPEN-API-KEY for this -> I suggest using GEMINI as it is free-to-use.

👨‍💻 Author
Abhishek Kudwa

GitHub: @Kudwa-Abhishek

🙏 Acknowledgments
Movie data provided by Gavin Lon (Youtube)
Icon generated from ChatGPT
Inspiration from modern streaming platforms, golang, langchain, open-source.
**BIG THANKS TO: GAVIN LON (Youtube) + FreeCodeCamp(Youtube) + GoLang(Documentation)**

📧 Contact
Email me at: ashu.kudwa@gmail.com
Linkedin: https://www.linkedin.com/in/abhishekrkudwa
