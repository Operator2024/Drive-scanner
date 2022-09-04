# Drive-scanner

## How it works

⚠️ The application will return either json or nothing. In some cases there may be panic

```
Usage of Drive-scanner:
  -V    This key allows you to get the current version
```

**Example #1**:

```bash
> docker build -t drive-scanner .
> docker run -ti --rm --privileged drive-scanner
```

```json
{
  "Drive-scanner": [
    {
      "/dev/sda": {....},
      "/dev/sdb": {....}
    }
  ]
}
```

## License

See the [LICENSE](LICENSE) file for license rights and limitations (MIT).
