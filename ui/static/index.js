const dataFile = {
    fileZIP: new FormData(),
    filePDF: null
}

const inputUploadFile = document.getElementById("input_upload")
inputUploadFile.addEventListener("change", (event) => {
    const file = event.target.files[0]

    if (file.type !== "application/zip") {
        return
    }

    if (file.size > 2147483648) {
        return
    }

    dataFile.fileZIP.set("zip", file)
})

const uploadFile = async () => {
    const response = await fetch("/api/upload", {
        method: "POST",
        body: dataFile.fileZIP
    })

    const json = await response.json()

    console.log(json)
}

const buttonSubmit = document.getElementById("submit_upload")
buttonSubmit.addEventListener("click", uploadFile)
