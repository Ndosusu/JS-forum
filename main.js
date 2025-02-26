import { createRegistrationForm } from "./scripts/registration.js";
import { createIndex } from "./scripts/index.js";
import { createProfile } from "./scripts/profile.js";
import { createContact } from "./scripts/contact.js";
import { createNavigation } from "./scripts/navigation.js";

function loadPage(page) {
    // Supprime uniquement le contenu généré par la page actuelle, pas la navbar
    const existingContent = document.getElementById("page-content");
    if (existingContent) {
        existingContent.remove();
    }

    // Vérifie si la navbar existe déjà, sinon la crée (sauf pour registration)
    const navigation = document.querySelector(".navigation");
    if (!navigation && page !== "registration") {
        createNavigation();
    } else if (navigation && page === "registration") {
        navigation.remove(); // Supprime la navbar si on est sur la page d'inscription
    }

    // Crée un conteneur pour la nouvelle page
    const newContent = document.createElement("div");
    newContent.id = "page-content";
    document.body.appendChild(newContent); // Ajoute le contenu sous la navbar

    // Charge la page demandée
    if (page === "registration") {
        createRegistrationForm(newContent);
    } else if (page === "index") {
        createIndex(newContent);
    } else if (page === "profile") {
        createProfile(newContent);
    } else if (page === "contact") {
        createContact(newContent);
    }
}

// Charge la page initiale
document.addEventListener("DOMContentLoaded", () => {
    loadPage("registration"); // Charge la page d'accueil par défaut
});

// Permet d'appeler loadPage() ailleurs
window.loadPage = loadPage;
