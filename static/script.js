function togglePassword(inputId, checkbox) {
  const input = document.getElementById(inputId);
  input.type = checkbox.checked ? "text" : "password";
}
