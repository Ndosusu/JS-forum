export function createIndex(parentElement) {
    const container = document.createElement("div");
    container.classList.add("container");

    const title = document.createElement("h2");
    title.textContent = "Accueil";
    container.appendChild(title);

    // New post part
    const newPostDiv = document.createElement("div");
    newPostDiv.id = "newPost-div";

    const newPostForm = document.createElement("form");
    newPostForm.id = "newPost-form";
    newPostForm.classList.add("newpostcontainer");

    const newPostLabel = document.createElement("p");
    newPostLabel.id = "newPost-label";
    newPostLabel.className = "";
    newPostLabel.innerHTML = "New Post";
    newPostForm.appendChild(newPostLabel);

    const newPostTitle = document.createElement("input");
    newPostTitle.id = "newPost-title";
    newPostTitle.name = "title";
    newPostTitle.placeholder = "Titre";
    newPostForm.appendChild(newPostTitle);

    const newPostContent = document.createElement("textarea");
    newPostContent.id = "newPost-content";
    newPostContent.name = "content";
    newPostContent.placeholder = "Contenu";
    newPostForm.appendChild(newPostContent);

    const newPostButton = document.createElement("button");
    newPostButton.textContent = "Jarvis, post that";
    newPostButton.onclick = async (event) => {
        event.preventDefault();

        const form = document.getElementById("newPost-form");
        const formData = new FormData(form);

        const data = {};
        formData.forEach((value, key) => (data[key] = value));

        const createPostResp = await fetch("/api/newPost", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        });

        if (createPostResp.ok) {
            alert("Post créé avec succès!");

            FetchPosts(postsContainer)

        } else {
            const error = await createPostResp.text();
            alert("Erreur lors de la création du post: " + error);
        }
    }
    newPostForm.appendChild(newPostButton);

    newPostDiv.appendChild(newPostForm);
    container.appendChild(newPostDiv);

    // Posts
    const postsContainer = document.createElement("section");
    postsContainer.id = "postsContainer";

    FetchPosts(postsContainer);

    container.appendChild(postsContainer);

    // Ajoute tout au conteneur parent (et non document.body)
    parentElement.appendChild(container);
}

function FetchPosts(container) {
    container.innerHTML = '';

    fetch("/api/fetchAllPosts")
        .then(response => response.json())
        .then(allPosts => {
            if (!allPosts) {
                const noPost = document.createElement("p");
                noPost.innerHTML = "No posts at the moment."
                container.appendChild(noPost);
            } else {
                allPosts.forEach(post => {
                    container.appendChild(NewPost(post));
                });
            }
        });
}

function NewPost(post) {
    const newPost = document.createElement("div");
    newPost.className = "post";

    newPost.innerHTML = `
        <div class="user-part">
            <img src="${post["profile_picture"]}" alt="Profile Picture" width="64" height="64" />
            <p>${post["username"]}</p>
        </div>
        <div class="content-part">
            <h3>${post["title"]}</h3>
            <p>${post["content"]}</p>
        </div>
        <div class="stats-part">
            <p>${post["likes"]} likes</p>
            <p>${post["dislikes"]} dislikes</p>
        </div>
        <p>Posted on ${new Date(post["created_at"]).toString()}</p>
    `;

    return newPost;
}