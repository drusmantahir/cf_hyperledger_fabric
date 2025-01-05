import subprocess
import hashlib
import os
import numpy as np
from pypuf.simulation import ArbiterPUF
from pypuf.io import random_inputs
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives import serialization

class ChaincodeOperations:
    def __init__(self, channel_name, chaincode_name, orderer, orderer_tls, peer0_org1, peer0_org2, peer0_org1_tls_cert, peer0_org2_tls_cert):
        self.channel_name = channel_name
        self.chaincode_name = chaincode_name
        self.orderer = orderer
        self.orderer_tls = orderer_tls
        self.peer0_org1 = peer0_org1
        self.peer0_org2 = peer0_org2
        self.peer0_org1_tls_cert = peer0_org1_tls_cert
        self.peer0_org2_tls_cert = peer0_org2_tls_cert

    def _execute_command(self, command):
        try:
            result = subprocess.run(command, shell=True, check=True, text=True, capture_output=True)
            return result.stdout
        except subprocess.CalledProcessError as e:
            return f"Error: {e.stderr}"

    def create_key(self, key):
        command = f"""
        peer chaincode invoke \
          --waitForEvent \
          -o {self.orderer} \
          --ordererTLSHostnameOverride orderer.example.com \
          --tls \
          --cafile {self.orderer_tls} \
          --peerAddresses {self.peer0_org1} \
          --tlsRootCertFiles {self.peer0_org1_tls_cert} \
          --peerAddresses {self.peer0_org2} \
          --tlsRootCertFiles {self.peer0_org2_tls_cert} \
          -C {self.channel_name} \
          -n {self.chaincode_name} \
          -c '{{"function":"CreateKey","Args":["{key}"]}}'
        """
        return self._execute_command(command)

    def create_key_with_image(self, key, image_url):
        command = f"""
        peer chaincode invoke \
          --waitForEvent \
          -o {self.orderer} \
          --ordererTLSHostnameOverride orderer.example.com \
          --tls \
          --cafile {self.orderer_tls} \
          --peerAddresses {self.peer0_org1} \
          --tlsRootCertFiles {self.peer0_org1_tls_cert} \
          --peerAddresses {self.peer0_org2} \
          --tlsRootCertFiles {self.peer0_org2_tls_cert} \
          -C {self.channel_name} \
          -n {self.chaincode_name} \
          -c '{{"function":"CreateKeyWithImage","Args":["{key}", "{image_url}"]}}'
        """
        return self._execute_command(command)

    def query_all_keys(self):
        command = f"""
        peer chaincode query \
          -o {self.orderer} \
          --ordererTLSHostnameOverride orderer.example.com \
          --tls \
          --cafile {self.orderer_tls} \
          --peerAddresses {self.peer0_org1} \
          --tlsRootCertFiles {self.peer0_org1_tls_cert} \
          -C {self.channel_name} \
          -n {self.chaincode_name} \
          -c '{{"function":"QueryAllKeys","Args":[]}}'
        """
        return self._execute_command(command)

    def read_key(self, key):
        command = f"""
        peer chaincode query \
          -o {self.orderer} \
          --ordererTLSHostnameOverride orderer.example.com \
          --tls \
          --cafile {self.orderer_tls} \
          --peerAddresses {self.peer0_org1} \
          --tlsRootCertFiles {self.peer0_org1_tls_cert} \
          -C {self.channel_name} \
          -n {self.chaincode_name} \
          -c '{{"function":"ReadKey","Args":["{key}"]}}'
        """
        return self._execute_command(command)

    def consume_key(self, key):
        command = f"""
        peer chaincode invoke \
          --waitForEvent \
          -o {self.orderer} \
          --ordererTLSHostnameOverride orderer.example.com \
          --tls \
          --cafile {self.orderer_tls} \
          --peerAddresses {self.peer0_org1} \
          --tlsRootCertFiles {self.peer0_org1_tls_cert} \
          --peerAddresses {self.peer0_org2} \
          --tlsRootCertFiles {self.peer0_org2_tls_cert} \
          -C {self.channel_name} \
          -n {self.chaincode_name} \
          -c '{{"function":"ConsumeKey","Args":["{key}"]}}'
        """
        return self._execute_command(command)

    
    def generate_puf_key(self):
        """
        Generate a PUF-based key and return the public key in hexadecimal format.
        """
        # Simulate PUF responses
        puf = ArbiterPUF(n=256, seed=45)
        fixed_challenge = np.array(random_inputs(n=256, N=1, seed=1))
        response = puf.eval(fixed_challenge)
        response_bits = (response > 0).astype(int).flatten()
        response_bytes = response_bits.tobytes()

        # Derive ECC private key
        private_key_int = int.from_bytes(hashlib.sha256(response_bytes).digest(), 'big')
        private_key = ec.derive_private_key(private_key_int, ec.SECP256R1(), default_backend())
        public_key = private_key.public_key()

        # Convert public key to hex
        public_key_hex = public_key.public_bytes(
            encoding=serialization.Encoding.X962,
            format=serialization.PublicFormat.UncompressedPoint
        ).hex()

        return public_key_hex


# Instantiate and Test
chaincode_ops = ChaincodeOperations(
    channel_name="unicalchannel",
    chaincode_name="keygen",
    orderer="localhost:7050",
    orderer_tls="/Users/drusmantahir/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem",
    peer0_org1="localhost:7051",
    peer0_org2="localhost:9051",
    peer0_org1_tls_cert="/Users/drusmantahir/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt",
    peer0_org2_tls_cert="/Users/drusmantahir/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"
)

#print(chaincode_ops.create_key("11223344556688"))
#print(chaincode_ops.create_key_with_image("11223344556699", "http://example.com/image.jpg"))
#print(chaincode_ops.query_all_keys())
print(chaincode_ops.read_key("11223344556699"))
#print(chaincode_ops.consume_key("11223344556688"))