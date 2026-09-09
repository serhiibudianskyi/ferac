<template>
  <div class="sp-acc">
    <div
      v-if="wallet"
      class="shadow-std acc-dd-btn flex items-center p-3 rounded-lg mr-3 hover:bg-gray-100 text-sm font-bold"
      :class="[state.accountDropdown ? 'active' : '']"
      style="display: flex; align-items: center"
      @click="state.accountDropdown = !state.accountDropdown"
    >
      <div class="flex items-center">
        <IgntProfileIcon :address="state.keplrParams?.bech32Address" />
        <span class="mx-2">
          {{ getAccName() }}
        </span>
      </div>
    </div>
    <IgntButton
      v-else
      aria-label="Connect wallet"
      type="primary"
      @click="state.connectWalletModal = true"
    >
      Connect wallet
    </IgntButton>
    <IgntAccDropdown
      v-if="state.accountDropdown && wallet"
      :wallet="wallet"
      :acc-name="getAccName()"
      :dev-wallets="walletStore.getDevWallets"
      @disconnect="disconnect"
      @select-wallet="selectWallet"
      @add-wallet="openWalletImport"
      @close="state.accountDropdown = false"
    />
    <IgntModal
      :visible="state.connectWalletModal"
      :close-icon="false"
      :cancel-button="false"
      :submit-button="false"
      style="text-align: center"
      @close="state.connectWalletModal = false"
      @submit="state.connectWalletModal = false"
    >
      <template #header>
        <div
          v-if="state.modalPage === 'connect'"
          class="flex items-center flex-col my-3"
        >
          <h3 v-if="isKeplrAvailable" class="text-2xl font-bold">
            Connect your wallet
          </h3>
          <h3 v-else>Connect your wallet</h3>
        </div>
        <div v-else-if="state.modalPage === 'connecting'">
          <div class="description-grey">Opening Keplr</div>
          <h3>Connecting</h3>
        </div>
        <div v-else-if="state.modalPage === 'mnemonic'">
          <h3>Import dev wallet</h3>
        </div>
        <div v-else-if="state.modalPage === 'error'">
          <IgntWarningIcon style="margin-bottom: 20px" />
          <h3>{{ state.errorMessage }}</h3>
        </div>
      </template>
      <template #body>
        <div class="max-w-xs text-center text-sm my-4 mx-auto">
          <div v-if="state.modalPage === 'connect'">
            <p v-if="isKeplrAvailable">
              Connect your Keplr wallet via the Keplr browser extension to use
              this app.
            </p>
            <p v-else>
              Connect with a local development mnemonic
            </p>
          </div>
          <div v-else-if="state.modalPage === 'connecting'">
            <div class="mt-8">
              <IgntSpinner />
            </div>
            <IgntButton
              aria-label="Cancel"
              type="secondary"
              style="margin-top: 3rem"
              @click="state.modalPage = 'connect'"
            >
              Cancel
            </IgntButton>
            <div class="external-link mt-8">Having trouble opening Keplr?</div>
          </div>
          <div v-else-if="state.modalPage === 'mnemonic'" class="text-left">
            <label class="text-xs text-gray-600" for="dev-wallet-name">
              Wallet name
            </label>
            <input
              id="dev-wallet-name"
              v-model="state.localWalletName"
              class="mt-1 mb-4 py-2 px-4 h-12 bg-gray-100 border-xs text-base leading-tight w-full rounded-xl outline-0"
              placeholder="Local Dev Wallet"
            />
            <label class="text-xs text-gray-600" for="dev-wallet-mnemonic">
              Mnemonic
            </label>
            <textarea
              id="dev-wallet-mnemonic"
              v-model="state.mnemonic"
              class="mt-1 py-2 px-4 min-h-28 bg-gray-100 border-xs text-base leading-tight w-full rounded-xl outline-0 resize-none"
              placeholder="Enter a local development mnemonic"
            />
            <p class="mt-3 text-xs text-gray-500 text-center">
              The mnemonic is used only in this browser session and is not saved.
            </p>
          </div>
          <div v-else-if="state.modalPage === 'error'" style="padding: 20px 0">
            <div class="external-link">
              <span>Keplr troubleshooting</span>
              <IgntExternalArrowIcon style="margin-left: 0.5rem" />
            </div>
          </div>
        </div>
      </template>
      <template #footer>
        <div v-if="state.modalPage === 'connect'" class="my-3">
          <div style="gap: 10px; display: flex; justify-content: center; flex-wrap: wrap">
            <IgntButton
              v-if="isKeplrAvailable"
              aria-label="Connect Keplr"
              type="primary"
              @click="tryToConnectToKeplr"
            >
              Connect Keplr
            </IgntButton>
            <IgntButton
              aria-label="Import local mnemonic"
              type="secondary"
              @click="state.modalPage = 'mnemonic'"
            >
              Import mnemonic
            </IgntButton>
          </div>
        </div>
        <div
          v-if="state.modalPage === 'mnemonic'"
          style="gap: 10px; display: flex; justify-content: center"
        >
          <IgntButton
            aria-label="Cancel mnemonic import"
            type="secondary"
            @click="state.modalPage = 'connect'"
          >
            Cancel
          </IgntButton>
          <IgntButton
            aria-label="Connect local wallet"
            type="primary"
            :disabled="state.mnemonic.trim().split(/\s+/).length < 12"
            @click="tryToConnectWithMnemonic"
          >
            Connect
          </IgntButton>
        </div>
        <div
          v-if="state.modalPage === 'error'"
          style="gap: 10px; display: flex; justify-content: center"
        >
          <IgntButton
            aria-label="Connect Keplr"
            type="secondary"
            @click="state.connectWalletModal = false"
          >
            Cancel
          </IgntButton>
          <IgntButton
            aria-label="Connect Keplr"
            type="primary"
            @click="state.modalPage = 'connect'"
          >
            Try again
          </IgntButton>
        </div>
      </template>
    </IgntModal>
  </div>
</template>

<script setup lang="ts">
import { IgntButton } from "@ignt/vue-library";
import { IgntExternalArrowIcon } from "@ignt/vue-library";
import { IgntKeplrIcon } from "@ignt/vue-library";
import { IgntModal } from "@ignt/vue-library";
import { IgntProfileIcon } from "@ignt/vue-library";
import { IgntSpinner } from "@ignt/vue-library";
import { IgntWarningIcon } from "@ignt/vue-library";
import { computed, onMounted, reactive, watch } from "vue";

import { useClient } from "@/composables/useClient";
import useCosmosBaseTendermintV1Beta1 from "@/composables/useCosmosBaseTendermintV1Beta1";
import useKeplr from "@/def-composables/useKeplr";
import { useWalletStore } from "@/stores/useWalletStore";

import IgntAccDropdown from "./IgntAccDropdown.vue";

export interface State {
  modalPage: string;
  connectWalletModal: boolean;
  accountDropdown: boolean;
  keplrParams: { name: string; bech32Address: string };
  localWalletName: string;
  mnemonic: string;
  errorMessage: string;
}

const initialState: State = {
  modalPage: "connect",
  connectWalletModal: false,
  accountDropdown: false,
  keplrParams: { name: "", bech32Address: "" },
  localWalletName: "Local Dev Wallet",
  mnemonic: "",
  errorMessage: "Wallet cannot connect",
};

// state
const state = reactive(initialState);

// composables
const { connectToKeplr, isKeplrAvailable, getKeplrAccParams } = useKeplr();

const client = useClient();
const walletStore = useWalletStore();
// methods
const wallet = computed(() => walletStore.getWallet);
const query = useCosmosBaseTendermintV1Beta1();
const nodeInfo = query.ServiceGetNodeInfo({});
const chainId = computed(
  () => nodeInfo.data?.value?.default_node_info?.network ?? ""
);
watch(
  () => chainId.value,
  async (newVal) => {
    if (newVal != "" && wallet.value?.name === "Keplr Integration") {
      const { name, bech32Address } = await getKeplrAccParams(newVal);
      state.keplrParams.name = name;
      state.keplrParams.bech32Address = bech32Address;
    }
  }
);

const tryToConnectToKeplr = async (): Promise<void> => {
  state.modalPage = "connecting";

  const onKeplrConnect = async () => {
    state.connectWalletModal = false;
    state.modalPage = "connect";
  };

  const onKeplrError = (): void => {
    state.modalPage = "error";
  };

  await connectToKeplr(onKeplrConnect, onKeplrError);
};
const tryToConnectWithMnemonic = async (): Promise<void> => {
  try {
    await walletStore.connectWithMnemonic(
      state.mnemonic,
      state.localWalletName.trim() || "Local Dev Wallet"
    );
    state.keplrParams.name = state.localWalletName.trim() || "Local Dev Wallet";
    state.keplrParams.bech32Address = walletStore.getAddress;
    state.mnemonic = "";
    state.connectWalletModal = false;
    state.modalPage = "connect";
  } catch (e) {
    console.error(e);
    state.errorMessage = "Mnemonic cannot be imported";
    state.modalPage = "error";
  }
};
const getAccName = (): string => {
  return wallet.value?.name ?? "";
};
const disconnect = async (): Promise<void> => {
  state.accountDropdown = false;
  await walletStore.signOut();
};
const openWalletImport = (): void => {
  state.accountDropdown = false;
  state.modalPage = "mnemonic";
  state.connectWalletModal = true;
};
const selectWallet = async (address: string): Promise<void> => {
  try {
    await walletStore.switchToMnemonicWallet(address);
    state.keplrParams.name = walletStore.getWallet?.name ?? "";
    state.keplrParams.bech32Address = walletStore.getAddress;
    state.accountDropdown = false;
  } catch (e) {
    console.error(e);
    state.errorMessage = "Wallet cannot be switched";
    state.accountDropdown = false;
    state.modalPage = "error";
    state.connectWalletModal = true;
  }
};

// check if already connected
onMounted(async () => {
  try {
    if (await walletStore.restoreMnemonicWallet()) {
      state.keplrParams.name = walletStore.getWallet?.name ?? "";
      state.keplrParams.bech32Address = walletStore.getAddress;
      return;
    }
  } catch (e) {
    console.warn("Saved local wallet could not be restored", e);
    walletStore.signOut();
  }

  if (client.signer) {
    try {
      await tryToConnectToKeplr();
    } catch (_e) {
      console.warn("Keplr not connected");
    }
  }
});
</script>
