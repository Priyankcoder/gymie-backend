const bcrypt = require('bcrypt');

const password = 'pri@expo';
const saltRounds = 10;

bcrypt.hash(password, saltRounds, (err, hash) => {
  if (err) {
    console.error('Error generating hash:', err);
    process.exit(1);
  }
  console.log('Bcrypt hash for password "pri@expo":');
  console.log(hash);
});
