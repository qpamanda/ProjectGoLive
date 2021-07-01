# ProjectGoLive - by Group 2 of GoSchool Run 4

#### **MEMBERS**
    - Ahmad Bahrudin
    - Amanda Soh Chiew Pheng
    - Huang Yanping
    - Tan Jun Jie

#### **INTRODUCTION**

Project Name: Giving Grace Portal

![Logo](/images/logo.jpg)

A platform for helping the needy, elderly, or charity organisations. An application that allows representatives for the needy (individuals or organisations) to sign up and put up requests for the needs required. The needy (individuals or organisation) will be our recipients. The requests to put up can be for a donation (monetary), donation (items) or errands. Example of an errand "XXX needs lunch/dinner for date/time". A helper who is also a user of the application can view the list of request and select a request to fulfill.


![tech](/images/Software Architecture.JPG)


![DMBS](/images/DBMS.JPG)

    - For every tables created we will include following fields for accountability purpose.
        - CreatedBy
        - Created_dt
        - LastModifiedBy
        - LastModified_dt   
 

![Stakeholders](/images/Stakeholders.JPG)


The following is the discussion and jobs delegations for our group thus far:

#### **APPLICATION FEATURES**

![Features](/images/Features.JPG)

###### 1)  USER LOGIN/REGISTRATION FEATURE (FOR REPRESENTATIVES - REQUESTERS/HELPERS) - AMANDA

    **Representatives Table**
    - RepID INT (Primary key , unique identifier for the account)
    - UserName
    - Password
    - FirstName
    - LastName
    - Email 
    - Password
    - FirstName
    - LastName
    - ContactNo
    - Organisation (optional)
    - CurrentLogin_dt
    - LastLogin_dt
    - CreatedBy
    - Created_dt
    - LastModifiedBy
    - LastModified_dt
    
    **RepMemberType Table**
    - RepID INT (Primary key , unique identifier for the account)
    - MemberTypeID (1 - Admin, 2 - Requester, 3 - Helper)
    - CreatedBy
    - Created_dt
    - LastModifiedBy
    - LastModified    
   

###### 2) RECIPIENT FEATURE (INDIVIDUAL/ORGANISATION WHO NEED HELP) - YANPING

    Recipent is the person or organisation needing help. Each Recipient is tied to a Representative
    
    **Recipients Table**
    - RecipientID INT (Primary key)
    - RepID
    - Name
    - Category (boolean - true for individual , false for organisation)
    - Profile
    - ContactNo
    - Email?
    - CreatedBy
    - Created_dt
    - LastModifiedBy
    - LastModified_dt


###### 3) REQUEST FEATURE (THE REQUESTS NEEDED BY THE RECIPIENT) - JUN JIE

    Requests made by individuals or organisations. To handle CRUD Request (Eg Add request . I need xxxx, update request status ) 
    
     **Requests Table**
     - RequestID (Primary key , unique identifier for the request)
     - RepID
     - CategoryID (1 - Item Donation, 2 - Errand)
     - RecipientID
     - RequestStatusCode (P - Pending, H - Being Handled, C - Completed)
     - RequestDetails
     - ToCompleteBy_dt
     - FulfilledAt (Location/Address - maybe not specific location but rather a region. Eg East , SouthEast area?)
     - CreatedBy
     - Created_dt
     - LastModifiedBy
     - LastModified_dt
      
      
###### 4) HELPERS FEATURE (PERSON HELPING THE REQUEST) - AMANDA
    
    The helper who has selected to fulfil the request(s).
    
    **Helpers Table**
    - RepID 
    - RequestID
    - CreatedBy
    - Created_dt
    - LastModifiedBy
    - LastModified_dt
    

###### 5) ADMIN MODULES (FOR SYSTEM SETUP) - AHMAD

    - CategoryID/CategoryName: 1 - Item Donation, 2 - Errand
    
    **Category Table**
    - CategoryID
    - Category
    - CreatedBy
    - Created_dt
    - LastModifiedBy
    - LastModified_dt
    
    
    - MemberTypeID/MemberType: 1 - Admin, 2 - Requester, 3 - Helper
    
    **MemberType Table**
    - MemberTypeID
    - MemberType
    - CreatedBy
    - Created_dt
    - LastModifiedBy
    - LastModified_dt


    - StatusCode/RequestStatus: P - Pending, H - Being Handled, C - Completed
    
    **RequestStatus Table**
    - StatusCode
    - Status
    - CreatedBy
    - Created_dt
    - LastModifiedBy
    - LastModified_dt


#### **BASIC FEATURES - COMPLETED**

    - [x] Connection to server via HTTPS/TLS (@port 5221) - Uses OpenSSL to create certificates
    - [x] Connection to database (using docker container @port 55005)
    - [x] Login/Logout
    - [x] Session Management
    - [x] Events Logging


#### **LIMITATIONS**    
The application cannot validate if the request is really fullfilled or not thus this is solely based on trust/honesty and helpers will liaise with representatives outside of the application to complete the request. The request status can then be updated into the system when it is completed.   


#### **FEATURES FOR IMPROVEMENTS**      

- [ ] Chat system ? (To facilitate the helpers and requesters) 
- [ ] Web Service for Helpers/Requesters to update request status?