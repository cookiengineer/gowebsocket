
(function() {

	let button   = document.querySelector("footer button#send");
	let main     = document.querySelector("main");
	let elements = {
		user:    document.querySelector("header input[name=\"user\"]"),
		room:    document.querySelector("header select[name=\"room\"]"),
		status:  document.querySelector("header span#status"),
		message: document.querySelector("footer input[name=\"message\"]"),
	};
	let socket = null;

	const render = (message) => {

		let div   = document.createElement("div");
		let span1 = document.createElement("span");
		let span2 = document.createElement("span");

		div.className = "message";

		span1.className   = "user";
		span1.textContent = message["user"] + ": ";

		span2.className   = "text";
		span2.textContent = message["text"];

		div.appendChild(span1);
		div.appendChild(span2);

		main.appendChild(div);
		main.scrollTop = main.scrollHeight;

	};

	const status = (state) => {
		elements["status"].textContent = state;
		elements["status"].setAttribute("data-status", state);
	};

	const connect = (user, room) => {

		if (socket !== null) {
			socket.close();
			socket = null;
		}

		let protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
		let url      = protocol + "//" + window.location.host + "/api/chat/" + encodeURIComponent(room);

		socket = new WebSocket(url);

		socket.onopen = () => {
			status("connected");
			button.removeAttribute("disabled");
		};

		socket.onmessage = (event) => {

			try {

				var message = JSON.parse(event.data);

				render(message);

			} catch (err) {
				console.error("Invalid message", err);
			}

		};

		socket.onclose = () => {
			status("disconnected");
			button.setAttribute("disabled", true);
			socket = null;
		};

		socket.onerror = () => {
			status("error");
			button.setAttribute("disabled", true);
			socket = null;
		};

	};

	const send = (user, message) => {

		if (socket.readyState === WebSocket.OPEN) {

			let payload = JSON.stringify({
				user: user,
				text: message,
			});

			socket.send(payload);

		}

	};

	const randomName = () => {

		let adjectives = [ "ambitious", "brave", "calm", "caring", "curious", "deceptive", "dedicated", "diligent", "empathetic", "emotional", "energetic", "funny", "honest", "humble", "kind", "loyal", "patient", "wise" ];
		let animals    = [ "armadillo", "cat", "dingo", "dog", "eagle", "hawk", "lion", "penguin", "tiger", "panda", "pangolin", "owl", "rabbit", "turtle", "wolf", "zebra" ];

		let tmp1 = adjectives[Math.floor(Math.random() * adjectives.length)];
		let tmp2 = animals[Math.floor(Math.random() * animals.length)];

		return tmp1 + "-" + tmp2;

	};

	if (elements["user"] !== null) {

		elements["user"].value = randomName();

		elements["user"].addEventListener("change", () => {

			if (elements["user"].checkValidity()) {
				connect(elements["user"].value, elements["room"].value);
			}

		});

	}

	if (elements["room"] !== null) {

		elements["room"].value = "#welcome";

		elements["room"].addEventListener("change", () => {

			if (elements["user"].checkValidity()) {
				connect(elements["user"].value, elements["room"].value);
			}

		});

	}

	if (button !== null) {

		button.addEventListener("click", () => {

			if (elements["user"].checkValidity()) {

				send(elements["user"].value, elements["message"].value);

				elements["message"].value = "";
				elements["message"].focus();

			}

		});

	}

	if (elements["message"] !== null) {

		elements["message"].addEventListener("keydown", (event) => {

			if (event.key === "Enter") {

				event.preventDefault();

				if (elements["user"].checkValidity()) {

					send(elements["user"].value, elements["message"].value);

					elements["message"].value = "";
					elements["message"].focus();

				}

			}

		});

	}

	console.log(elements);

	if (elements["user"] !== null && elements["room"] !== null) {
		connect(elements["user"].value, elements["room"].value);
	}

})();

