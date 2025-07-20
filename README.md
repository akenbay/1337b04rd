# 1337b04rd
## Structure:
### Package: `domain/`  
- **Purpose**: Stores all the needed structures and describes which functions should be in the interface.


### Package: `services/`  
- **Purpose**: Holds structure of the services which holds the interfaces with functions. Operate core business logic using those functions from the interfaces.
- **Depends On**: `domain/` for simple structures of requests, users, comments, and posts, also for interfaces

