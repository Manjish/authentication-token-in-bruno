const fetch = require("node-fetch");

//required variables from environment
const baseUrl = bru.getEnvVar("base_url");
const basicAuthUsername = bru.getEnvVar("basic_auth_username");
const basicAuthPassword = bru.getEnvVar("basic_auth_password");
const adminEmail = bru.getEnvVar("admin_email");
const adminPassword = bru.getEnvVar("admin_password");
const studentEmail = bru.getEnvVar("student_email");
const studentPassword = bru.getEnvVar("student_password");

const getToken = async (userType) => {
  // if requesting userType is not valid
  if (["admin", "student"].indexOf(userType) === -1) {
    console.log("Invalid user type");
    throw new Error("Invalid user type");
  }

  let tokenType = "";
  let tokenSetAtType = "";
  let loginEmail = "";
  let loginPassword = "";

  // setting variables based on userType
  switch (userType) {
    case "admin":
      tokenType = "adminToken";
      tokenSetAtType = "adminTokenSetAt";
      loginEmail = adminEmail;
      loginPassword = adminPassword;
      break;
    case "student":
      tokenType = "studentToken";
      tokenSetAtType = "studentTokenSetAt";
      loginEmail = studentEmail;
      loginPassword = studentPassword;
      break;
  }

  // check if token is already set and is valid
  const token = bru.getVar(tokenType);

  if (token) {
    const tokenSetAt = bru.getVar(tokenSetAtType);
    const currentTime = Date.now();
    const timeDifference = currentTime - tokenSetAt;
    const timeDifferenceInSeconds = timeDifference / 1000;

    // token is valid for 1 hour
    if (timeDifferenceInSeconds < 3600) {
      console.log(`${userType} token is still valid`);
      return;
    }
  }

  // if token is not set or is expired, get a new token
  const response = await fetch(`${baseUrl}/auth/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Basic ${Buffer.from(
        `${basicAuthUsername}:${basicAuthPassword}`
      ).toString("base64")}`,
    },
    body: JSON.stringify({
      email: loginEmail,
      password: loginPassword,
    }),
  });

  const responseData = await response.json();

  // if response is ok, set the token and time when it was set
  if (response.ok) {
    bru.setVar(tokenType, responseData.data.IdToken);
    bru.setVar(tokenSetAtType, Date.now());
  }
};

module.exports = getToken;
