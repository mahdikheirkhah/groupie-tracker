# Groupie-Tracker

## Objectives

- Receiv a given API and manipulate the data contained in it in order to create a website displaying the information.

- Using the given API, that is made up of four parts:

- **Artists: Contains information about some bands and artists like their name(s), image, in which year they began their activity, the date of their first album and the members.**

- **Locations: Contains their last and/or upcoming concert locations.**

- **Dates: Contains their last and/or upcoming concert dates.**

- **Relation: Links the data of all the other parts, artists, dates and locations.**

- Given all this, we built a user-friendly website which can display the band's info through several data visualizations.

- This project also focuses on the creation and visualization of events/actions.

---

## Project Structure

```
    groupie-tracker/
    ├── server.go               # Main application server
    ├── go.mod
    ├── utils/
    │   ├── errors.go
    │   ├── pagehandler.go
    │   ├── readfromapi.go
    │   ├── saferender.go
    │   └── structs.go
    ├── templates/
    │   ├── index.html
    │   ├── soncers.html
    │   ├── badRequest.html
    │   ├── internalServer.html
    │   └── notFound.html
    └── static/
        ├── index.html
        └── css/
            ├── styles.css
            └── index.html
```

## Usage

### How to Run

1. Clone this repository to your local machine:
   ```bash
   git clone https://01.gritlab.ax/git/mkheirkh/groupie-tracker.git
   ```
2. Navigate to the project directory:
   ```bash
   cd groupie-tracker
   ```
3. Ensure you have Go installed (version 1.19+ recommended).

4. Run the application:

   ```bash
   go run .

   ```

5. Open your browser and navigate to http://localhost:8080

---

## Instructions

1. Browse thru the page to select a band or use the search bar.
2. Select either "See Details" at search bar or "More Information" under the artist picture.
3. Check thru the artist page.

---

- Routes:

  - GET /:
    - Renders the homepage using the index.html template.
    - Serves static CSS files for styling.

- Error Handling:
  - Custom HTML templates (badRequest.html, notFound.html, internalServer.html) are served for HTTP error codes.

---

## Authors

- Mohammad mahdi Kheirkhah
- Fatemeh Kheirkhah
- Toft Diederichs

---
