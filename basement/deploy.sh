#!/usr/bin/env bash

# Stop if a command fails
set -e

if [ $# -ne 2 ]; then
  echo "Please provide <product-version> and <enigmas key>"
  exit 1
fi

PRODUCT_VERSION=$1
ENIGMAS_KEY="$2";

# DEFAULT VALUES IN DEPLOY VARIABLES #---------------------------

LOG_LEVEL="${LOG_LEVEL:-INFO}"
POSTGRES_INIT_DB="${POSTGRES_INIT_DB:-SOFT}"

# ---------------------------------------------------------------

# BASEMENT REPO UPDATED CHECK #----------------------------------
cd /opt/basement
git fetch
if [ $(git rev-parse HEAD) != $(git rev-parse @{u}) ]; then
  git merge origin/master
  echo "Basement branch was not updated. Please run again."
  exit 1
fi
# ---------------------------------------------------------------


# GET YARD AND BARCKARD VERSIONS --------------------------------

cd /opt/basement
YARD_VERSION=$(awk "\$5 == \"$PRODUCT_VERSION\" { print \$1 }" VERSIONS)
BACKYARD_VERSION=$(awk "\$5 == \"$PRODUCT_VERSION\" { print \$3 }" VERSIONS)

if [ -z "$YARD_VERSION" ] || [ -z "$BACKYARD_VERSION" ]; then
  echo "There was a problem fetching yard (`$YARD_VERSION`) or backyard (`$BACKYARD_VERSION`) versions"
  exit 1
fi
echo "Versions to be deployed are YARD=$YARD_VERSION and BACKYARD=$BACKYARD_VERSION."

# ---------------------------------------------------------------


idempotent_checkout() {
  local branch_or_tag=$1
  git fetch --all --tags
  if [ -n "$(git branch --list "$branch_or_tag")" ]; then
    git checkout "$branch_or_tag"
    git merge "origin/$branch_or_tag"
  elif [ -n $(git tag --list "$branch_or_tag") ]; then
    git checkout "$branch_or_tag"
  else
    git checkout --track "origin/$branch_or_tag"
  fi
}



# CLONING REPOS -------------------------------------------------

cd /opt
if [ ! -d /opt/yard ]; then
    echo "Cloning Yard";
    git clone git@github.com:hospedate/yard.git;
    echo "Cloning Yard completed.";
fi
if [ ! -d /opt/backyard ]; then
    echo "Cloning Backyard"
    git clone git@github.com:hospedate/backyard.git
    echo "Cloning Backyard completed."
fi

# ---------------------------------------------------------------


# CLEANING DOCKER IMGS ------------------------------------------
echo "Cleaning up old and dangling docker images"
export IMAGE_LIMIT=4

output=$(docker images | grep backyard | tail -n +"$IMAGE_LIMIT" | awk '{print $2}')
if [ -n "$output" ]; then
  docker image rm -f $(docker images | grep backyard | tail -n +"$IMAGE_LIMIT" | awk '{print $3}')
fi

output=$(docker images | grep yard | grep -v backyard | tail -n +"$IMAGE_LIMIT" | awk '{print $2}')
if [ -n "$output" ]; then
  docker image rm -f $(docker images | grep yard | grep -v backyard | tail -n +"$IMAGE_LIMIT" | awk '{print $3}')
fi

docker image prune -f

echo "Cleaning up old docker images completed."

# ---------------------------------------------------------------

echo "Building Yard"
cd /opt/yard
idempotent_checkout $YARD_VERSION
docker build --progress plain . -t yard:$YARD_VERSION --build-arg PRODUCT_VERSION=$PRODUCT_VERSION
echo "Building Yard completed."

echo "Building Backyard"
cd /opt/backyard
idempotent_checkout $BACKYARD_VERSION
docker build --progress plain . -t backyard:$BACKYARD_VERSION
echo "Building Backyard completed."

echo "Unlocking DB enigma"
cd /opt/basement/enigmas
gpg --no-symkey-cache --batch --passphrase "$ENIGMAS_KEY" db_enigma.sh.gpg
source db_enigma.sh # load var $DB_ENIGMA
rm -f db_enigma.sh

echo "Unlocking encryption key enigma"
cd /opt/basement/enigmas
gpg --no-symkey-cache --batch --passphrase "$ENIGMAS_KEY" encryption_key.sh.gpg
source encryption_key.sh # load var $ENCRYPTION_KEY
rm -f encryption_key.sh

echo "Unlocking SMTP enigma"
cd /opt/basement/enigmas
gpg --no-symkey-cache --batch --passphrase "$ENIGMAS_KEY" smtp_password.sh.gpg
source smtp_password.sh # load var $SMTP_PASSWORD
rm -f smtp_password.sh

cd /opt/basement
BACKYARD_VERSION=$BACKYARD_VERSION YARD_VERSION=$YARD_VERSION DB_ENIGMA="$DB_ENIGMA" POSTGRES_INIT_DB=$POSTGRES_INIT_DB PRODUCT_VERSION=$PRODUCT_VERSION LOG_LEVEL=$LOG_LEVEL SMTP_PASSWORD=$SMTP_PASSWORD ENCRYPTION_KEY=$ENCRYPTION_KEY docker-compose up --force-recreate -d
