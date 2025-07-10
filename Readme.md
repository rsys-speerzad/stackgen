# Coding Exercise Details

## Objective

In a geographically distributed team, it is very hard to find common time to meet that works for everyone. Your objective is to build an API, using which we will solve this problem for any event.

The organizer creates an event with a brief title “Brainstorming meeting” and provides. N slots (eg 12 Jan 2025, 2 - 4PM EST, 14 Jan 2025 6-9 PM EST etc.) also provide estimated time required for the meeting eg. 1 hr.

All participants also provide their availability in the similar format. 

The system recommends the time slots that work for all. If there is no such time slot found, then it recommends time slots that work for the most number of people (also provides a list for whom it does not work).

## Core Functionality

- Implement the REST API in Golang.

- Support creating, updating, and deleting events

- Support creating, updating, and deleting preferred time slots by each user.

- Endpoint that shows the possible time slots for the event.

## Expectations

### Must Haves:

- Automated tests for the written code.

- The application must be deployable on cloud infrastructure.

- Sticking to REST conventions

### Good to Have:

- Design for horizontal scalability.

- OpenAPI spec to document REST APIs.

- Containerization of the application.

- Provide Infrastructure as Code (IaC) (e.g., Terraform, Helm) to deploy.