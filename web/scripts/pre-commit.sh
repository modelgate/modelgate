#!/bin/sh

# 获取当前分支
branch_name=$(git rev-parse --abbrev-ref HEAD)

if [ "$branch_name" != "master" ]; then
  echo "🟡 On '$branch_name' branch, skipping typecheck and lint-staged."
  exit 0
fi

echo "✅ Running typecheck and lint-staged on branch '$branch_name'..."

pnpm typecheck && pnpm lint-staged

# 如果上面任何一步失败，就 fail 掉 commit
if [ $? -ne 0 ]; then
  echo "❌ Pre-commit checks failed."
  exit 1
fi

echo "🎉 Pre-commit checks passed."
exit 0
