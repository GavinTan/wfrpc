# compile for version
make
if [ $? -ne 0 ]; then
    echo "make error"
    exit 1
fi

wfrpc_version=`./bin/wfrpc --version`
echo "build version: $wfrpc_version"

# cross_compiles
make -f ./Makefile.cross-compiles

rm -rf ./release/packages
mkdir -p ./release/packages

os_all='linux windows darwin freebsd'
arch_all='386 amd64 arm arm64 mips64 mips64le mips mipsle'

cd ./release

for os in $os_all; do
    for arch in $arch_all; do
        wfrpc_dir_name="wfrpc_${wfrpc_version}_${os}_${arch}"
        wfrpc_path="./packages/wfrpc_${wfrpc_version}_${os}_${arch}"

        if [ "x${os}" = x"windows" ]; then
            if [ ! -f "./wfrpc_${os}_${arch}.exe" ]; then
                continue
            fi

            mkdir ${wfrpc_path}
            mv ./wfrpc_${os}_${arch}.exe ${wfrpc_path}/wfrpc.exe
        else
            if [ ! -f "./wfrpc_${os}_${arch}" ]; then
                continue
            fi

            mkdir ${wfrpc_path}
            mv ./wfrpc_${os}_${arch} ${wfrpc_path}/wfrpc
        fi

        # packages
        cd ./packages
        if [ "x${os}" = x"windows" ]; then
            zip -rq ${wfrpc_dir_name}.zip ${wfrpc_dir_name}
        else
            tar -zcf ${wfrpc_dir_name}.tar.gz ${wfrpc_dir_name}
        fi  
        cd ..
        rm -rf ${wfrpc_path}
    done
done

cd -
