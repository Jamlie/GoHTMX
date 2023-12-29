const css = `
	cursor: pointer;
`;

class Navbar extends HTMLElement {
	constructor() {
		super();
	}

	connectedCallback() {
		this.innerHTML = `
			<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
				<div class="container">
					<div class="navbar-nav">
						<a style="${css}" class="navbar-brand" hx-get="/home" hx-swap="outerHTML" hx-replace-url="true" hx-target="body">Home</a>
						<a style="${css}" class="navbar-brand" hx-get="/about" hx-swap="outerHTML" hx-replace-url="true" hx-target="body">About</a>
					</div>
					<div class="ml-auto">
						<a style="${css}" class="navbar-brand" hx-post="/api/logout" hx-swap="outerHTML" hx-replace-url="true" hx-target="body">Logout</a>
					</div>
				</div>
			</nav>
		`;
	}
}

customElements.define('nav-bar', Navbar);
