package UniswapV2Staker

const StakingRewardsBin = "0x608060405260006004556000600555624f1a0060065534801561002157600080fd5b506040516113073803806113078339818101604052606081101561004457600080fd5b508051602082015160409092015160018055600280546001600160a01b039485166001600160a01b031991821617909155600380549285169282169290921790915560008054939092169216919091179055611262806100a56000396000f3fe608060405234801561001057600080fd5b50600436106101415760003560e01c80637b0a47ee116100b8578063cd3daf9d1161007c578063cd3daf9d146102ad578063d1af0c7d146102b5578063df136d65146102bd578063e9fad8ee146102c5578063ebe2b12b146102cd578063ecd9ba82146102d557610141565b80637b0a47ee1461025257806380faa57d1461025a5780638b87634714610262578063a694fc3a14610288578063c8f33c91146102a557610141565b8063386a95251161010a578063386a9525146101d35780633c6b16ab146101db5780633d18b912146101f85780633fc6df6e1461020057806370a082311461022457806372f702f31461024a57610141565b80628cc262146101465780630700037d1461017e57806318160ddd146101a45780631c1f78eb146101ac5780632e1a7d4d146101b4575b600080fd5b61016c6004803603602081101561015c57600080fd5b50356001600160a01b031661030d565b60408051918252519081900360200190f35b61016c6004803603602081101561019457600080fd5b50356001600160a01b03166103a3565b61016c6103b5565b61016c6103bc565b6101d1600480360360208110156101ca57600080fd5b50356103da565b005b61016c610569565b6101d1600480360360208110156101f157600080fd5b503561056f565b6101d16107c0565b6102086108e4565b604080516001600160a01b039092168252519081900360200190f35b61016c6004803603602081101561023a57600080fd5b50356001600160a01b03166108f3565b61020861090e565b61016c61091d565b61016c610923565b61016c6004803603602081101561027857600080fd5b50356001600160a01b0316610931565b6101d16004803603602081101561029e57600080fd5b5035610943565b61016c610acc565b61016c610ad2565b610208610b2c565b61016c610b3b565b6101d1610b41565b61016c610b64565b6101d1600480360360a08110156102eb57600080fd5b5080359060208101359060ff6040820135169060608101359060800135610b6a565b6001600160a01b0381166000908152600a6020908152604080832054600990925282205461039d919061039190670de0b6b3a7640000906103859061036090610354610ad2565b9063ffffffff610d8c16565b6001600160a01b0388166000908152600c60205260409020549063ffffffff610de916565b9063ffffffff610e4916565b9063ffffffff610eb316565b92915050565b600a6020526000908152604090205481565b600b545b90565b60006103d5600654600554610de990919063ffffffff16565b905090565b60018054810190819055336103ed610ad2565b6008556103f8610923565b6007556001600160a01b0381161561043f576104138161030d565b6001600160a01b0382166000908152600a60209081526040808320939093556008546009909152919020555b60008311610488576040805162461bcd60e51b8152602060048201526011602482015270043616e6e6f74207769746864726177203607c1b604482015290519081900360640190fd5b600b5461049b908463ffffffff610d8c16565b600b55336000908152600c60205260409020546104be908463ffffffff610d8c16565b336000818152600c60205260409020919091556003546104ea916001600160a01b039091169085610f0d565b60408051848152905133917f7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5919081900360200190a2506001548114610565576040805162461bcd60e51b815260206004820152601f6024820152600080516020611199833981519152604482015290519081900360640190fd5b5050565b60065481565b6000546001600160a01b031633146105b85760405162461bcd60e51b815260040180806020018281038252602a8152602001806111da602a913960400191505060405180910390fd5b60006105c2610ad2565b6008556105cd610923565b6007556001600160a01b03811615610614576105e88161030d565b6001600160a01b0382166000908152600a60209081526040808320939093556008546009909152919020555b60045442106106395760065461063190839063ffffffff610e4916565b600555610688565b60045460009061064f904263ffffffff610d8c16565b9050600061066860055483610de990919063ffffffff16565b60065490915061068290610385868463ffffffff610eb316565b60055550505b600254604080516370a0823160e01b815230600482015290516000926001600160a01b0316916370a08231916024808301926020929190829003018186803b1580156106d357600080fd5b505afa1580156106e7573d6000803e3d6000fd5b505050506040513d60208110156106fd57600080fd5b505160065490915061071690829063ffffffff610e4916565b600554111561076c576040805162461bcd60e51b815260206004820152601860248201527f50726f76696465642072657761726420746f6f20686967680000000000000000604482015290519081900360640190fd5b426007819055600654610785919063ffffffff610eb316565b6004556040805184815290517fde88a922e0d3b88b24e9623efeb464919c6bf9f66857a65e2bfcf2ce87a9433d9181900360200190a1505050565b60018054810190819055336107d3610ad2565b6008556107de610923565b6007556001600160a01b03811615610825576107f98161030d565b6001600160a01b0382166000908152600a60209081526040808320939093556008546009909152919020555b336000908152600a6020526040902054801561089b57336000818152600a6020526040812055600254610864916001600160a01b039091169083610f0d565b60408051828152905133917fe2403640ba68fed3a2f88b7557551d1993f84b99bb10ff833f0cf8db0c5e0486919081900360200190a25b505060015481146108e1576040805162461bcd60e51b815260206004820152601f6024820152600080516020611199833981519152604482015290519081900360640190fd5b50565b6000546001600160a01b031681565b6001600160a01b03166000908152600c602052604090205490565b6003546001600160a01b031681565b60055481565b60006103d542600454610f64565b60096020526000908152604090205481565b6001805481019081905533610956610ad2565b600855610961610923565b6007556001600160a01b038116156109a85761097c8161030d565b6001600160a01b0382166000908152600a60209081526040808320939093556008546009909152919020555b600083116109ee576040805162461bcd60e51b815260206004820152600e60248201526d043616e6e6f74207374616b6520360941b604482015290519081900360640190fd5b600b54610a01908463ffffffff610eb316565b600b55336000908152600c6020526040902054610a24908463ffffffff610eb316565b336000818152600c6020526040902091909155600354610a51916001600160a01b03909116903086610f7a565b60408051848152905133917f9e71bc8eea02a63969f509818f2dafb9254532904319f9dbda79b67bd34a5f3d919081900360200190a2506001548114610565576040805162461bcd60e51b815260206004820152601f6024820152600080516020611199833981519152604482015290519081900360640190fd5b60075481565b6000600b5460001415610ae857506008546103b9565b6103d5610b1d600b54610385670de0b6b3a7640000610b11600554610b11600754610354610923565b9063ffffffff610de916565b6008549063ffffffff610eb316565b6002546001600160a01b031681565b60085481565b336000908152600c6020526040902054610b5a906103da565b610b626107c0565b565b60045481565b6001805481019081905533610b7d610ad2565b600855610b88610923565b6007556001600160a01b03811615610bcf57610ba38161030d565b6001600160a01b0382166000908152600a60209081526040808320939093556008546009909152919020555b60008711610c15576040805162461bcd60e51b815260206004820152600e60248201526d043616e6e6f74207374616b6520360941b604482015290519081900360640190fd5b600b54610c28908863ffffffff610eb316565b600b55336000908152600c6020526040902054610c4b908863ffffffff610eb316565b336000818152600c602052604080822093909355600354835163d505accf60e01b81526004810193909352306024840152604483018b9052606483018a905260ff8916608484015260a4830188905260c4830187905292516001600160a01b039093169263d505accf9260e480820193929182900301818387803b158015610cd257600080fd5b505af1158015610ce6573d6000803e3d6000fd5b5050600354610d0992506001600160a01b0316905033308a63ffffffff610f7a16565b60408051888152905133917f9e71bc8eea02a63969f509818f2dafb9254532904319f9dbda79b67bd34a5f3d919081900360200190a2506001548114610d84576040805162461bcd60e51b815260206004820152601f6024820152600080516020611199833981519152604482015290519081900360640190fd5b505050505050565b600082821115610de3576040805162461bcd60e51b815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b600082610df85750600061039d565b82820282848281610e0557fe5b0414610e425760405162461bcd60e51b81526004018080602001828103825260218152602001806111b96021913960400191505060405180910390fd5b9392505050565b6000808211610e9f576040805162461bcd60e51b815260206004820152601a60248201527f536166654d6174683a206469766973696f6e206279207a65726f000000000000604482015290519081900360640190fd5b6000828481610eaa57fe5b04949350505050565b600082820183811015610e42576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b604080516001600160a01b038416602482015260448082018490528251808303909101815260649091019091526020810180516001600160e01b031663a9059cbb60e01b179052610f5f908490610fda565b505050565b6000818310610f735781610e42565b5090919050565b604080516001600160a01b0385811660248301528416604482015260648082018490528251808303909101815260849091019091526020810180516001600160e01b03166323b872dd60e01b179052610fd4908590610fda565b50505050565b610fec826001600160a01b0316611192565b61103d576040805162461bcd60e51b815260206004820152601f60248201527f5361666545524332303a2063616c6c20746f206e6f6e2d636f6e747261637400604482015290519081900360640190fd5b60006060836001600160a01b0316836040518082805190602001908083835b6020831061107b5780518252601f19909201916020918201910161105c565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d80600081146110dd576040519150601f19603f3d011682016040523d82523d6000602084013e6110e2565b606091505b509150915081611139576040805162461bcd60e51b815260206004820181905260248201527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564604482015290519081900360640190fd5b805115610fd45780806020019051602081101561115557600080fd5b5051610fd45760405162461bcd60e51b815260040180806020018281038252602a815260200180611204602a913960400191505060405180910390fd5b3b15159056fe5265656e7472616e637947756172643a207265656e7472616e742063616c6c00536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f7743616c6c6572206973206e6f742052657761726473446973747269627574696f6e20636f6e74726163745361666545524332303a204552433230206f7065726174696f6e20646964206e6f742073756363656564a265627a7a72315820f2b812934500f8f3e0d662a34a9600d72498f3a503e73560bf849a66533976b864736f6c63430005100032"
