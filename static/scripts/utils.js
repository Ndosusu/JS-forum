import { createSignContainer } from "./registration.js";
import { createIndex } from "./index.js";
import { createProfile } from "./profile.js";
import { createContact } from "./contact.js";
import { createNavigation } from "./navigation.js";
import { createListOfUsers } from "./listOfUsers.js";
import { createCategories } from "./categories.js";


export function LoadPage(page) {
    // Supprime uniquement le contenu généré par la page actuelle, pas la navbar
    const content = document.getElementById("page-content");
    if (content) {
        content.innerHTML = "";
    }

    // Vérifie si la navbar existe déjà, sinon la crée (sauf pour registration)
    const navigation = document.querySelector(".navigation");
    const categories = document.querySelector(".categories");
    const listUsers = document.querySelector(".listContainer")
    if (!navigation && page !== "registration") {
        createNavigation();
        createListOfUsers();
    } else if (navigation && page === "registration") {
        navigation.remove(); // Supprime la navbar si on est sur la page d'inscription
        listUsers.remove();
    }

    if (!categories && page === "index") {
        createCategories();
    } else if (categories && page !== "index") {
        categories.remove();
    }

    // Charge la page demandée
    switch (page) {
        case "contact":
            createContact(content);
            break;
        
        case "index":
            createIndex(content);
            break;

        case "profile":
            createProfile(content);
            break;
        
        case "registration":
            createSignContainer(content);
            break;
    }
}

// Remove quotes contained in a string
function removeQuotes(uuid) {
    return uuid.replace(/"/g, '');
}

// Parse user info contained inside cookie
export function GetUserInfoFromCookie() {
    if (document.cookie.length === 0) {
        return null;
    }

    let userInfo = {
        uuid: "",
        username: "",
        email: "",
        role: ""
    };

    if (document.cookie.substring(0, 10) === 'UserLogged') {
        const parts = document.cookie.substring(11).split('|');

        if (parts.length >= 4) {
            userInfo = {
                uuid: removeQuotes(parts[0]),   // UUID
                username: parts[1],             // Nom d'utilisateur
                email: parts[2],                // Email
                role: parts[3]                  // Rôle
            };

            if (parts.length === 5) {
                userInfo.profileImageURL = removeQuotes(parts[4]);  // URL de l'image de profil
            }
        }
    }

    return userInfo;
}

// Is the cookie returned by getUserInfoFromCookie valid?
export function IsUserInfoValid() {
    var userInfo = GetUserInfoFromCookie();
    var uuidRegex = /([0-9a-f]{8})-([0-9a-f]{4})-([0-9a-f]{4})-([0-9a-f]{4})-([0-9a-f]{12})/;
    var emailRegex = /([0-9A-Za-z]+[\.-]*)+@([0-9A-Za-z]+-*)+.(com|org|fr)/;

    // No cookie
    if (!userInfo) {
        console.log("No cookie");
        return false;
    }

    // Is UUID valid?
    if (!uuidRegex.test(userInfo.uuid)) {
        console.log("UUID not valid");
        return false;
    }

    // Is email valid?
    if (!emailRegex.test(userInfo.email)) {
        console.log("Email not valid");
        return false;
    }

    // Is role valid?
    switch (userInfo.role) {
        case "":
        case "user":
        case "mod":
        case "admin":
            break;
        default:
            return false;
    }

    return true;
}