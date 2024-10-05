“* **Database Connection** * 
[x] Successfully connects to MongoDB using the MongoDB Go driver. * **User Registration** * [x] Allows users to register with a username, display name, and profile picture. * 
[x] Hashes passwords using bcrypt before storing them in the database. * 
[x] Inserts new user documents into the `users` collection. * **User Login** * 
[x] Authenticates users by checking the username and validating the hashed password. * 
[x] Generates a JWT upon successful login, which includes the username and expiration time. * **JWT Generation and Validation** * 
[x] Creates JWTs signed with a secret key (your-256-bit-secret). * 
[x] Validates JWTs to ensure they were signed correctly and haven't expired. * **CORS Support** * 
[x] Implements CORS to allow requests from the specified frontend origin ([http://localhost:3000](http://localhost:3000)). * **WebSocket Support** * 
[x] Sets up a WebSocket endpoint (`/ws`) for real-time communication (requires a valid JWT for access). * **Error Handling** * 
[x] Handles errors during user registration, login, password verification, and JWT generation. * 
[x] Provides appropriate HTTP status codes and error messages for various scenarios (e.g., user not found, invalid password, etc.). * **Modular Code Organization** * 
[x] Uses separate functions for registration, login, password hashing, and JWT management, promoting code reusability and clarity. * **Graceful MongoDB Disconnection** * 
[x] Ensures that the MongoDB client is disconnected when the application is terminated.”


