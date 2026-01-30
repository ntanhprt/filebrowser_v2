@echo off
set CGO_ENABLED=0

echo === Build frontend ===
pushd frontend
call pnpm install
call pnpm build
popd

echo Frontend build complete.
pause

