package main

import (
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)
var upgrader = websocket.Upgrader{
       EnableCompression: true,
       CheckOrigin:       func(r *http.Request) bool { return true },
}

func DefineAllRoutes(router *gin.Engine){
    router.GET("/ws", wsHandler)
    router.GET("/", home)
}

func home(c *gin.Context){
    c.String(http.StatusOK, "Welcome to the home page")
}
func wsHandler(c *gin.Context){
    user := c.Query("user")
    if user == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user query parameter"})
        return
    }
    //c.String(http.StatusOK, "WebSocket handler placeholder",user)

//     // Upgrade HTTP -> WebSocket
   conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Println("Failed to upgrade:", err)
        //c.String(http.StatusInternalServerError, "WebSocket upgrade failed: %v", err)
        return
    }

    // ✅ Correct way to print connection info
    //c.String(http.StatusOK, "WebSocket connection established: %v", conn)
    defer conn.Close()
    log.Println("Client connected:", user)

    // Echo loop`
    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("Read error:", err)
            break
        }
        log.Printf("[%s] says: %s", user, string(message))

        // Send back to the client (echo)
        if err := conn.WriteMessage(messageType, message); err != nil {
            log.Println("Write error:", err)
            break
        }
    }
    log.Println("Client disconnected:", user)
}
func main() {
    router := gin.Default()
    DefineAllRoutes(router)
    router.Run(":8080")
}
