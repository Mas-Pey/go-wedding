const token = localStorage.getItem('token')
if (!token) {
    window.location.href = 'login.html'
}

async function fetchProtected(url) {
    const token = localStorage.getItem('token')
    const headers = {
        'Authorization': 'Bearer ' + token,
        'Content-Type': 'application/json'
    }

    try {
        const response = await fetch(url, {headers: headers})
        if (response.status === 401) {
            alert("Silakan login ulang. Expired session, please relogin.")
            localStorage.removeItem('token')
            window.location.href = 'login.html'
            return null
        }
        return response
    } catch (error) {
        alert("Failed to connect server.")
        return null
    }
}

function Logout() {
    const logoutBtn = document.getElementById('logoutBtn')
    if (logoutBtn) {
        logoutBtn.addEventListener('click', function(e) {
            e.preventDefault()
            localStorage.removeItem('token')
            window.location.href = 'login.html'
        })
    }
}