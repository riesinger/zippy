document.addEventListener("DOMContentLoaded", function() {
    console.log("--- LOADED ---");

    var simplemde = new SimpleMDE({
                element: document.getElementById("editor"),
                autoDownloadFontAwesome: false,
                toolbar: false
            });

    var saveButton = document.body.querySelector("button#save");
    var titleInput = document.body.querySelector("input#post-title");
    var draftCheckbox = document.body.querySelector("input#draft-checkbox");

    saveButton.addEventListener("click", function() {
        var postContent = simplemde.value();
        var postData = JSON.stringify({
            title: titleInput.value,
            isDraft: draftCheckbox.checked,
            body: postContent
        });
        fetch("/api/saveArticle", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: postData
        })
            .then(function(response) {
                return response.json();
            })
            .then(function(response) {

            })
            .catch(function(err) {
                console.error(err);
            });
    });

});
