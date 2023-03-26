profile="$DEPLOY_PROFILE"
branch="$DEPLOY_BRANCH"

if [[ -z $DEPLOY_PROFILE ]]; then
  echo "DEPLOY_PROFILE env variable not set"
fi

case "$profile" in
  "prod")
    branch="main"
    ;;

  "test")
    branch="develop"
    ;;

  *)
    echo "Invalid DEPLOY_PROFILE value. Must be 'prod' or 'test'"
    exit 1
    ;;
esac

echo "Profile: $profile; Branch: $branch"

cd "~/diverse/Diverse-Back-$profile"
git checkout "$branch"
git pull
sudo docker-compose -f docker-compose.yml --profile "$profile" stop
sudo docker-compose -f docker-compose.yml --profile "$profile" pull
sudo docker-compose -f docker-compose.yml --profile "$profile" up -d --build
