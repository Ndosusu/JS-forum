import { createRegistrationForm } from "./scripts/registration.js";
import { createIndex } from "./scripts/index.js";
import { createProfile } from "./scripts/profile.js";
import { createContact } from "./scripts/contact.js";

function loadPage(page) {
    // Supprime uniquement le contenu généré par la page actuelle, pas la navbar
    const existingContent = document.getElementById("page-content");
    if (existingContent) {
        existingContent.remove();
    }

    // Crée un conteneur pour la nouvelle page
    const newContent = document.createElement("div");
    newContent.id = "page-content";
    document.body.appendChild(newContent); // Ajoute le contenu sous la navbar

    //retire la navbar sur la page d'inscription/connexion
    const navigation = document.querySelector(".navigation");
    if (navigation) {
        if (page === "registration") {
            navigation.style.display = "none";
        } else {
            navigation.style.display = "block";
        }
    }
    // Charge la page demandée
    // ajouter le nom de la page dans la navbar
    if (page === "registration") {
        createRegistrationForm(newContent);
    } else if (page === "index") {
        createIndex(newContent);
    } else if (page === "profile") {
        createProfile(newContent);
    } else if (page === "contact") {
        createContact(newContent)
    }
}


document.addEventListener("DOMContentLoaded", () => {
            loadPage("registration");
        });


// Permet d'appeler loadPage() ailleurs
window.loadPage = loadPage;
