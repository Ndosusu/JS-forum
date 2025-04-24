export function createSignContainer(parentElement) {
    const container = document.createElement("div");
    container.classList.add("container");

    const title = document.createElement("h2");
    title.textContent = "Inscription / Connexion";
    container.appendChild(title);

    container.appendChild(createLoginForm());
    container.appendChild(createRegistrationForm());

    // Bouton pour aller à l'accueil
    const goToIndexButton = document.createElement("button");
    goToIndexButton.textContent = "debug home";
    goToIndexButton.addEventListener("click", () => {
        LoadPage("index");
    });

    container.appendChild(goToIndexButton);

    // Ajoute tout au conteneur parent
    parentElement.appendChild(container);
}

function createLoginForm() {
    const loginForm = document.createElement("form");
    loginForm.id = "login-form";

    /*
        <input type="email" id="log-email" placeholder="Email" required>
        <input type="password" id="log-password" placeholder="Password" required>
        <button onclick="sendLoginForm">Se connecter</button>
    */

    const emailInput = document.createElement('input');
    emailInput.type = "email";
    emailInput.id = "log-" + emailInput.type;
    emailInput.placeholder = "Email";
    emailInput.required = true;
    loginForm.appendChild(emailInput);

    const passwordInput = document.createElement("input");
    passwordInput.type = "password";
    passwordInput.id = "log-" + passwordInput.type;
    passwordInput.placeholder = "Password";
    passwordInput.required = true;
    loginForm.appendChild(passwordInput);

    const submitBtn = document.createElement("button");
    submitBtn.innerHTML = "Se connecter";
    submitBtn.onclick = async (event) => {
        // Fonction pour envoyer le formulaire de connexion au serveur (en JSON)
        event.preventDefault();
        event.target.disabled = true;

        const email = document.getElementById("log-email").value;
        const password = document.getElementById("log-password").value;

        const data = {
            email: email,
            password: password
        }

        try {
            const response = await fetch("/api/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(data)
            });

            if (response.ok) {
                window.location.href = "/";
            } else {
                const error = await response.json();
                alert("Erreur lors du login: " + error.message);
            }
        } catch (error) {
            console.error("Erreur lors du login", error);
        }
    };
    loginForm.appendChild(submitBtn);

    return loginForm;
}

function createRegistrationForm() {
    const registerForm = document.createElement("form");
    registerForm.id = "register-form";

    /*
        <input type="text" id="reg-username" name="username" placeholder="Username" required />
        <input type="email" id="reg-email" name="email" placeholder="Email" required />
        <input type="password" id="reg-password" name="password" placeholder="Password" required />
        <input type="text" id="reg-firstName" name="firstName" placeholder="First Name" required />
        <input type="text" id="reg-lastName" name="lastName" placeholder="Last Name" required />
        <input type="date" id="reg-birthDate" name="birthDate" placeholder="Birth Date" required />
        <label for="gender">Gender:</label>
        <select id="reg-gender" name="gender">
            <option value="">Please select one…</option>
            <option value="female">Female</option>
            <option value="male">Male</option>
            <option value="non-binary">Non-Binary</option>
            <option value="other">Other</option>
            <option value="Prefer not to answer">Prefer not to Answer</option>
        </select>
        <button onclick="sendRegisterForm">S'inscrire</button>
    */

    const usernameInput = document.createElement("input");
    usernameInput.type = "text";
    usernameInput.name = "username";
    usernameInput.id = "reg-" + usernameInput.name;
    usernameInput.placeholder = "Username";
    usernameInput.required = true;
    registerForm.appendChild(usernameInput);

    const emailInput = document.createElement("input");
    emailInput.type = "email";
    emailInput.name = "email";
    emailInput.id = "reg-" + emailInput.name;
    emailInput.placeholder = "Email";
    emailInput.required = true;
    registerForm.appendChild(emailInput);

    const passwordInput = document.createElement("input");
    passwordInput.type = "password";
    passwordInput.name = "password";
    passwordInput.id = "reg-" + passwordInput.name;
    passwordInput.placeholder = "Password";
    passwordInput.required = true;
    registerForm.appendChild(passwordInput);

    const firstnameInput = document.createElement("input");
    firstnameInput.type = "text";
    firstnameInput.name = "firstName";
    firstnameInput.id = "reg-" + firstnameInput.name;
    firstnameInput.placeholder = "First Name";
    firstnameInput.required = true;
    registerForm.appendChild(firstnameInput);

    const lastnameInput = document.createElement("input");
    lastnameInput.type = "text";
    lastnameInput.name = "lastName";
    lastnameInput.id = "reg-" + lastnameInput.name;
    lastnameInput.placeholder = "Last Name";
    lastnameInput.required = true;
    registerForm.appendChild(lastnameInput);

    const birthdateInput = document.createElement("input");
    birthdateInput.type = "date";
    birthdateInput.name = "birthDate";
    birthdateInput.id = "reg-" + birthdateInput.name;
    birthdateInput.required = true;
    registerForm.appendChild(birthdateInput);

    // Gender
    const genderInput = document.createElement("select");
    genderInput.name = "gender";
    genderInput.id = "reg-" + genderInput.name;
    genderInput.innerHTML = `
        <option value="">Please select one…</option>
        <option value="female">Female</option>
        <option value="male">Male</option>
        <option value="non-binary">Non-Binary</option>
        <option value="other">Other</option>
        <option value="Prefer not to answer">Prefer not to Answer</option>
    `;
    registerForm.appendChild(genderInput);

    const submitBtn = document.createElement("button");
    submitBtn.innerHTML = "S'inscrire";
    submitBtn.onclick = async (event) => {
        // Fonction pour envoyer le formulaire d'inscription au serveur (en JSON)
        event.preventDefault();
        event.target.disabled = true;

        const username = document.getElementById("reg-username").value;
        const password = document.getElementById("reg-password").value;
        const email = document.getElementById('reg-email').value;
        const firstName = document.getElementById("reg-firstName").value;
        const lastName = document.getElementById("reg-lastName").value;
        const birthDate = document.getElementById("reg-birthDate").value;
        const gender = document.getElementById("reg-gender").value;

        // Prepare data object to be sent
        const data = {
            username: username,
            password: password,
            email: email,
            first_name: firstName,
            last_name: lastName,
            birth_date: birthDate + "T00:00:00Z",
            gender: gender
        };

        try {
            const response = await fetch("/api/register", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(data)  // Send the data as JSON
            });

            if (response.ok) {
                window.location.href = "/";
            } else {
                const error = await response.json();
                alert("Erreur lors de l'inscription : " + error.message);
            }
        } catch (error) {
            console.error("Erreur lors de l'inscription", error.message);
        }
    };
    registerForm.appendChild(submitBtn);

    return registerForm;
}