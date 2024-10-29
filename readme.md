# Ethereum & Arbitrum & Orbit Integration CLI Toolkit(Deprecate v1.0.0)
본 Repo는 개발자들이 Ethereum(Layer1) & Arbitrum(Layer2) & Orbit(Layer3) 솔루션을 원활하게 통합할 수 있도록 설계된 CLI 명령어와 라이브러리, 그리고 On-Chain 데이터 파싱 기능을 제공합니다.

This repository provides developers with powerful tools and libraries designed to seamlessly integrate Ethereum's Layer 1, Arbitrum's Layer 2, and Orbit's Layer 3 solutions.

## Quick Start
```bash
$ git submodule update --init

$ cp config.exmaple.yml config.yml

$ make build
```

## How to use
```bash
$ orbit-cli
```
## How to add account
### Create new account
```bash
$ orbit-cli account new

? Enter the password [for skip <ENTER>] :  
```

### Import account from private key
```bash
$ orbit-cli account import

? Enter the private key:  

? Enter the password [for skip <ENTER>] :  
```
| If set a password, need to enter it if a signature is required for all step. but, if no set password, will be skip enter password

## How to custom provider
Enter Custom Provider Url on config.yml
```yml
providers:
  local:
    - http://localhost:8545 # Layer 1
    - http://localhost:8547 # Layer 2
    - http://localhost:3347 # Layer 2
  sepolia: # chain tag
    - < Enter Sepolia Provider URL >
    - < Enter Sepolia Provider URL >
    - < Enter Sepolia Provider URL >
```


## Requirement
| Name | Version |
|------|---------|
| Go   | ^1.21.0 |
| Node | ^18     |
| Make | Latest  |
