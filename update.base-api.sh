remote_url=$(git remote get-url origin)
if [ $remote_url == "git@github.com:go-uniform/base-service.git" ] || [ $remote_url == "https://github.com/go-uniform/base-service.git" ] || [ $remote_url == "git@github.com:go-uniform/base-api.git" ] || [ $remote_url == "https://github.com/go-uniform/base-api.git" ]
then
  echo "The base-api repository can't be updated from within itself!"
  exit
fi
template_remote_url=$(git remote get-url template)
if [ $template_remote_url == "" ]
then
  git remote add template git@github.com:go-uniform/base-api.git
  template_remote_url=$(git remote get-url template)
fi

if [ $template_remote_url != "git@github.com:go-uniform/base-api.git" ] && [ $template_remote_url != "https://github.com/go-uniform/base-api.git" ]
then
  echo "The template remote url is incorrectly set, unable to update successfully!"
  exit
fi

git fetch template main
git merge template/main --allow-unrelated-histories
