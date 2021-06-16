# ProjectGoLive by Group 2

Members:
- Ahmad Bahrudin
- Amanda Soh Chiew Pheng
- Huang Yanping
- Tan Jun Jie

INTRODUCTION

A platform for helping the needy, elderly, or charity organisations. An application that allows representatives for the needy (individuals or organisations) to sign up and put up requests for the needs required. The needy (individuals or organisation) will be our recipients. The requests to put up can be for a donation (monetary), donation (items) or errands. Example of an errand "XXX needs lunch/dinner for date/time". A helper who is also a user of the application can view the list of request and select a request to fulfill.

Constraints: The application cannot validate if the request is really fullfilled or not thus this is solely based on trust/honesty and helpers will liaise with representatives outside of the application to complete the request. The request status can then be updated into the system when it is completed.

The following is the discussion and jobs delegations for our group thus far:

APPLICATION FEATURES

1)  ACCOUNT REGISTRATION FEATURE (FOR REPRESENTATIVES/HELPERS) - AMANDA

    Fields Required:
    
    RepresentativeID (Primary key , unique identifier for the account)
    Email 
    Password
    FirstName
    LastName
    PhoneNo
    Oranisation (optional)
    
    MemberType: Requestor, Helper?
    ..
    ..
    ..
    
    CRUD representative (To allow user manage their group. Kind of like adding a person into the Group in the trace together for family safe entry thingy)
    representative Fields:
    representative id (Primary key)
    Managed by (email, Foreign key)
    name ( the 1 receiving help)
    phone no.
    NRIC ???


2) RECIPIENT FEATURE (INDIVIDUAL/ORGANISATION NEED HELP) - YANPING

    Recipent is the person or organisation needing help
    
    Recipient to tie to a Representative
    RepresentativeID
    RecipentID (Individual) / RecipentID (Organisation) - Same tables or different tables?
    Recipient Info:
      ...
      ...
      ...


3) REQUEST FEATURE (THE REQUESTED NEEDED BY THE RECIPIENT) - JUN JIE

    CRUD Request (Eg Add request . I need xxxx, update request status ) 
    
      Fields Required:
        RequestID (Primary key , unique identifier for the request)
        RepresentativeID - Hosted by
        CategoryID    
        IndividualID --> RecipientID
        OrganisationID --> RecipientID
        ContactNo
        RequestStatus: Pending, Being Handled, Completed
        Request Details
          RequestDescription
          Date/time (Mainly for errands ?)
          Location/Address ?? (maybe not specific location but rather a region. Eg East , SouthEast area)
      
      
4) HELPER FEATURE (PERSON HELPING THE REQUEST) - AMANDA


5) ADMIN MODULE (FOR SYSTEM SETUP) - AHMAD

    CategoryID/CategoryName: 1 Donation (Monetary), 2 Donation (Physical Items), 3 Errands

    MemberTypeID/MemberType: Requestor, Helper

    StatusCode/RequestStatus: Pending, Being Handled, Completed



FEATURES - COMPLETED

6) Connection to server via https
7) Connection to database
8) Login/Logout
9) Session Management
10) Events Logging
11) 


FEATURES - ON HOLD      

7) Chat system ? (To facilitate the helpers and requestors) 


ADHOC    

- For every tables created we will include following fields for tracking purpose.
  CreatedBy
  Created_DT
  ModifiedBy
  Modified_DT   
    
    
    
      
    
    
    
    
