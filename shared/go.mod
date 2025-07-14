module github.com/VibeTeam/fitness-tracker-backend/shared

go 1.24

require (
    github.com/gin-gonic/gin v1.10.1
    github.com/VibeTeam/fitness-tracker-backend/user v0.0.0-00010101000000-000000000000
)

replace github.com/VibeTeam/fitness-tracker-backend/user => ../user
