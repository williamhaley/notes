<script lang="ts">
  import { logIn } from "./stores/auth";

  let form: HTMLFormElement;

  function handleSubmit(event: Event) {
    event.preventDefault();
    const formData = new FormData(form);

    const username = formData.get('username');
    const password = formData.get('password');

    if (!username || !password) {
      throw new Error('Please provide both username and password');
    }

    logIn(String(username), String(password));
  }
</script>

<form bind:this={form} class="container" on:submit={handleSubmit}>
  <h1>Log In</h1>

  <div class="form-group">
    <label for="email">Email address</label>
    <input type="email" class="form-control" id="email" name="username" placeholder="Enter email" required />
  </div>
  <div class="form-group">
    <label for="password">Password</label>
    <input type="password" class="form-control" id="password" name="password" placeholder="Password" required />
  </div>

  <button type="submit" class="btn btn-primary">Submit</button>
</form>
