profile=''

print_usage() {
  printf "Usage: deploy
  -p      Profile. May be prod or test.\n"
}

while getopts 'p:' flag; do
  case "${flag}" in
    p) profile="${OPTARG}" ;;
    *) print_usage
       exit 1 ;;
  esac
done

deploy() {
  local branch="$1"
  local profile="$2"

  if [[ -z $branch ]]; then
    echo "branch argument are empty"
  fi

  if [[ -z $profile ]]; then
    echo "profile argument are empty"
  fi

  echo $branch $profile

  cd "~/diverse/Diverse-Back-$profile"
  git checkout "$branch"
  git pull
  sudo docker-compose -f docker-compose.yml --profile "$profile" stop
  sudo docker-compose -f docker-compose.yml --profile "$profile" pull
  sudo docker-compose -f docker-compose.yml --profile "$profile" up -d --build
}

case $profile in
  "prod") deploy 'main' $profile ;;
  "test") deploy 'develop' $profile ;;
  *) print_usage
     exit 1 ;;
esac
