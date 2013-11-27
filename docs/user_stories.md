# Deliverable 3

Do deliverable 3 stuff first.  Then do the extra for extra bonus marks.

## From the prof

This deliverable askes use to implement the following use cases:

* Register Staff
* Staff Log in
* Staff Log out
* Register patient
* Consult Patient file
* Update Patient file
* Admit patient

They're described as :

* _Register Staff_ :
> Primary Actor: Staff Member (Nurse, Doctor, Medical Director, Personnel > Officer)
> Linked Requirements: F10
>
> Goal: A Staff Member wants to register to the HMS in order to identify
> himself to the system in order to use its functions
>
> Precondition: The HMS is ON
>
> Postcondition: The Staff Member is registered
>
> Main Scenario
>
> 1. Staff Member selects to register
>
> 2. HMS asks for Staff Member information
>
> 3. Staff Member provides all the required information
>
> 4. HMS displays an acknowledgment message
>
> Alternatives
>
> 3. a. Incomplete information is provided
>
> 3. a. 1. HMS displays an incomplete information error message
>
> 3.b. User not found in the system
>
> 3. b. 1. HMS displays an invalid user error message
>

* _Staff Log in_ :

> Primary Actor: Staff Member
>
> Linked Requirements: F12
>
> Goal: A Staff Member wants to identify himself to the system in order to use
> its functions
>
> Precondition: The HMS is ON AND Staff member is not logged in
>
> Postcondition: The Staff Member is logged in, the Staff Nember's operation
> choice menu is being displayed with the right priviledges
>
> Main Scenario
>
> 1. Staff Member enters login information
>
> 2. HMS checks Staff Member authorization status
>
> 3. HMS displays Staff Member operation choices
>
> Alternatives
>
> 2. a. Staff Member is not authorized to use system
>
> 2. a. 1. HMS displays unauthorized access error message
>


* _Staff Log out_ :

> Primary Actor: Staff Member
>
> Linked Requirements: F12
>
> Precondition: The Staff Member is logged in Postcondition: Staff Member is
> not > logged in Main Scenario
> 1. Staff Member selects to logs out
>
> 2. HMS displays log out acknowledgment
>

* _Register patient_ :
> Primary Actor: Medical Staff Member
>
> Linked Requirements: F1
>
> Goal: A Medical Staff Member wants to register a new Patient to the
> hospital. > The Patient is issued a Patient's identifier to be used at the hospital
> Precondition: The Medical Staff Member is logged in
>
> Postcondition: The Patient is registered and assigned an identification
> number
> Main Scenario
>
> 1. Medical Staff Member asks for new Patient registration
>
> 2. HMS asks for Patient's information
>
> 3. Medical Staff Member enters requested information
>
> 4. HMS issues a new identification number for Patient
>
> 5. HMS registers Patient
>
> Alternatives
>
> 3. a. Information is incomplete
>
> 3. a. 1. HMS displays incomplete information error message
>

* _Consult Patient file_ :

> Primary Actor: Staff Member
>
> Linked Requirements: F11
>
Goal: A Staff Member wants to visualize a Patient file. Precondition: The > Staff Member is logged in
> Postcondition: none
>
> Main Scenario
>
> 1. Staff Member asks for viewing Patient registration
>
> 2. HMS asks for Patient identification number
>
> 3. Medical staff member enters Patient identification number 4. HMS shows
> the > Patient registration information Alternatives
> 3. a. Identification number is incorrect
>
> 3. a. 1. HMS displays Patient not found error message
>

* _Update Patient file_ :

> Primary Actor: Medical Staff Member
>
> Linked Requirements: F2
>
> Goal: A Medical Staff Member wants to update Patient's information. >
> Precondition: The Medical Staff Member is logged in
> Postcondition: The Patient's information is updated
>
> Main Scenario
>
> 1. Include use case Consult Patient File
>
> 2. Medical Staff Member modifies information at will
>
> 3. Medical Staff Member resubmit
>
> 3. HMS updates Patient registration information
>
> Alternatives
>
> 3. a. Medical Staff Member do not have enough priviledge to modify 3. a. 1.
> > HMS displays modification not allowed error message

* _Admit patient_ :

> Primary Actor: Charge Nurse
>
> Linked Requirements: F4
>
> Precondition: The Charge Nurse is logged in
>
> Postcondition: A Patient is admitted to a division
>
> Main Scenario
>
> 1. Include use case Consult File
>
> 2. Charge Nurse chooses to admit the patient
>
> 3. HMS asks for room and bed number
>
> 4. Charge Nurse enters room and bed number or alternativelly browse through
> a list of available rooms and beds in the ward and makes a selection
> 5. HMS asks for remaining admission information
>
> 6. Charge Nurse enters requested information
>
> 7. HMS admits patient
>
> Alternatives
>
> 3. a. Division is complete
>
> 3. a. 1. HMS notifies Charge Nurse that her division is complete
>
> 3. a. 2. HMS gives possibiliy to request an admission for Patient as in use
> case Request Patient Admission

## User stories

* Register Staff

As a anonymous person, I want to create a staff account so that I can log in.

* Staff Log in

As a logged out staff, I want to log in so that I can access the staff features.

* Staff Log out

As a logged in staff, I want to log out so that I can remove access to my account on this browser.

* Register patient

As a logged in medical staff, I want to register a new Patient so that the patient file can be handled by the system.

* Consult Patient file

As a logged in medical staff, I want to consult a Patient file so that I can update the Patient file.

* Update Patient file

As a logged in medical staff, I want to update a Patient file so that I can make it relevant information.

* Admit patient

As a logged in Charge Nurse, I want to admit a Patient to my division so that staff in my division may treat him.

# Extras

## From the prof

All the use cases, minus those implemented above:

* Visualize division
* Request Patient Admission
* Admit patient from request list
* Discharge Patient
* Prescribe Medication


* _Visualize division_:

> Primary Actor: Charge Nurse
>
> Linked Requirements: F3
>
> Precondition: The Charge Nurse is logged in
>
> Postcondition: HMS is displaying information about a division
>
> Main Scenario
>
> 1. Charge Nurse chooses to visualize division
>
> 2. HMS asks for division identifier
>
3. Charge Nurse enters division identifier or alternatively browse to select > division 4. HMS displays information about division
> Alternatives
>
> 3. a. Wrong division identifier
>
> 3. a. 1. HMS displays division not found error message
>

* _Request Patient Admission_:

> Primary Actor: Charge Nurse
>
> Linked Requirements: F4
>
> Precondition: The Charge Nurse is logged in
>
> Postcondition: A Patient is admitted to a division
>
> Main Scenario
>
> 1. Include use case Consult File
>
> 2. Charge Nurse chooses to admit the patient
>
> 3. HMS asks for room and bed number
>
> 4. Charge Nurse enters room and bed number or alternativelly browse through
> a list of available rooms and beds in the ward and makes a selection
> 5. HMS asks for remaining admission information
>
> 6. Charge Nurse enters requested information
>
> 7. HMS admits patient
>
> Alternatives
>
> 3. a. Division is complete
>
> 3. a. 1. HMS notifies Charge Nurse that her division is complete
>
> 3. a. 2. HMS gives possibiliy to request an admission for Patient as in use
> case Request Patient Admission

* _Admit patient from request list_:

> Primary Actor: Charge Nurse
>
> Linked Requirements: F6
>
> Precondition: The Charge Nurse is logged in AND Patient is in request list
> Postcondition: Patients in request list admitted
> Main Scenario
>
> 1. Charge Nurse browses through list of Patients in request list and select
> one 2. HMS displays selected Patient registration
> 3. perform admission as in use case Admit Patient (steps 2 - )
>
> Alternatives
>
> 3. a. Patient can not be admitted
>
> 3. a. 1. Charge Nurse denies Patient admission
>
> 3. a. 2. HMS sends notification to Charge Nurse who requested admission
>

* _Discharge Patient_:

> Primary Actor: Charge Nurse
>
> Linked Requirements: F7
>
> Precondition: The Charge Nurse is logged in
>
> Postcondition: Patient is no longer admitted, Patient bed added to the
> division availabilities, Patient discharge information issued and copies slated to be sent to Patient external doctor, Patient invoice issued
> Main Scenario
>
> 1. Include use case Consult File
>
> 2. Charge Nurse selects to discharge Patient
>
> 3. HMS updates bed availability
>
> 4. HMS prints discharge information
>
> 5. HMS prints Patient invoice
>

* _Prescribe Medication_:

> Primary Actor: Doctor
>
> Linked Requirements: F8
>
> Precondition: The Doctor is logged
>
> Postcondition: Prescription has been added to the Patient file Main Scenario
>
> 1. Include use case Consult File
>
> 2. Doctor asks to add a prescription to the selected patient
>
> 3. HMS asks for required information
>
> 4. Doctor enters required information
>
> 5. HMS records prescription in Patient's file
>
> 2. a. Selected Patient is not one of the Doctor's
>
> 2. a. 1. HMS notifies Doctor that Patient is not hers
>
> 4. a. Incorrect information entered
>
> 4. a. 1. HMS display incorrect prescription error message
>

## User Stories

Write some user stories yo

* As a `{{role}}`, I want to `{{goal}}` so that `{{end result}}`.
