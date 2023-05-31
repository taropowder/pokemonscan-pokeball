mkdir -p build/pokeball/bin
mkdir -p build/pokeball/attachment
gitDescribe=$(git log --pretty=oneline --abbrev-commit -1)
go build  -ldflags "-X 'main.gitDescribe=${gitDescribe}'" -buildvcs=auto  -o ./build/pokeball/bin/pokeball ./src/cmd
cp script/pokeball.service build/pokeball/attachment
cp script/pokeball.sh build/pokeball/bin
chmod +x build/pokeball/bin/*
tar -czf  build/pokeball.tar.gz -C build/pokeball/ .