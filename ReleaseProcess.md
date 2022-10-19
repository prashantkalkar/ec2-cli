1. Merge all changes to main
2. Ensure main is upto date 
3. Ensure there are no local changes

4. Create tag

(Update following command by incrementing versions)
```shell
git tag v0.1
git push origin v0.1
```

5. Build binaries for platforms

```shell
make dist
```
6. Create release on Github UI

7. Upload assets (ec2-cli-darwin-amd64, ec2-cli-linux-amd64, ec2-cli-darwin-arm64) 


