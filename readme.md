# Ethereum & Arbitrum & Orbit Integration CLI Toolkit
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

## Custom Provider
Enter Custom Provider Url on config.yml
```yml
providers:
  local:
    ethereum: http://localhost:8545
    arbitrum: http://localhost:8547
    orbit: http://localhost:3347
  sepolia: # chain tag
    ethereum: < Enter Sepolia Provider URL >
    arbitrum: < Enter Sepolia Provider URL >
    orbit: < Enter Sepolia Provider URL >
```


## Requirement
| Name | Version |
|------|---------|
| Go   | ^1.21.0 |
| Node | ^18     |
| Make | Latest  |
