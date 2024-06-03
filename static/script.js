const socket = new WebSocket('ws://' + window.location.host + '/ws');

socket.onmessage = function(event) {
    const data = JSON.parse(event.data);
    if (data.error) {
        document.getElementById('output').innerText += 'Error: ' + data.error + '\n';
    } else if (data.status) {
        displayVMStatus(data.status);
    } else if (data.output) {
        document.getElementById('output').innerText += data.output + '\n';
    }
};

socket.onclose = function(event) {
    document.getElementById('output').innerText += 'Connection closed.\n';
};

document.getElementById('refresh').addEventListener('click', function() {
    sendCommand('status', '');
});

function sendCommand(command, arg) {
    const message = {
        command: command,
        arg: arg
    };
    document.getElementById('output').innerText = '';
    socket.send(JSON.stringify(message));
}

function displayVMStatus(statuses) {
    const tableBody = document.getElementById('vm-table-body');
    tableBody.innerHTML = '';
    statuses.forEach(status => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${status.Name}</td>
            <td>${status.State}</td>
            <td>${status.Port}</td>
            <td>${getActions(status.State, status.Name)}</td>
        `;
        tableBody.appendChild(row);
    });
}

function getActions(state, name) {
    let actions = '';
    if (state === 'running') {
        actions += `<button onclick="sendCommand('stop', '${name}')">Stop</button>`;
    } else if (state === 'poweroff') {
        actions += `<button onclick="sendCommand('start', '${name}')">Start</button>`;
    } else if (state === 'not_created') {
        actions += `<button onclick="sendCommand('start', '${name}')">Start</button>`;
    }
    actions += `<button onclick="sendCommand('reload', '${name}')">Reload</button>`;
    actions += `<button onclick="sendCommand('reboot', '${name}')">Reboot</button>`;
    actions += `<button onclick="sendCommand('remove', '${name}')">Remove</button>`;
    return actions;
}

window.onload = function() {
    fetch('/ipinfo')
        .then(response => response.json())
        .then(data => {
            document.getElementById('public-ip').innerText = data.publicIP;
            document.getElementById('private-ip').innerText = data.privateIP;
        });

    sendCommand('status', '');
};

