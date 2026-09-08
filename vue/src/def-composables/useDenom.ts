import BigNumber from "bignumber.js";
import { computed, ref } from "vue";

import useIbcApplicationsTransferV1 from "@/composables/useIbcApplicationsTransferV1";

const useDenomInstances = {} as Record<
  string,
  ReturnType<typeof useDenomInstance>
>;
const FERAC_BASE_DENOM = "uferac";
const FERAC_DISPLAY_DENOM = "FERAC";

export const formatDenomAmount = (amount: string, denom: string): string => {
  if (denom !== FERAC_BASE_DENOM) {
    return amount;
  }

  const parsedAmount = new BigNumber(amount || 0);
  if (!parsedAmount.isFinite()) {
    return "0.000000";
  }

  return parsedAmount.dividedBy(1_000_000).toFormat(6);
};
const traceToPath = (trace: { port_id: string; channel_id: string }[] | undefined) => {
  if (!trace || trace.length === 0) {
    return "";
  }else{
    return trace.map((t) => `${t.port_id}/${t.channel_id}`).join("/");
  }
};
const useDenomInstance = (denom: string) => {
  const isIBC = denom.indexOf("ibc/") == 0;
  const hash = denom.split("/")[1];
  const { QueryDenom } = useIbcApplicationsTransferV1();
  const denomTrace = QueryDenom(hash, { enabled: ref(isIBC) }).data;
  const normalized = computed(() => {
    if (isIBC) {
      return denomTrace.value?.denom?.base?.toUpperCase() ?? "";
    } else {
      if (denom === FERAC_BASE_DENOM) {
        return FERAC_DISPLAY_DENOM;
      }
      return denom.toUpperCase();
    }
  });
  const path = computed(() => {
    if (isIBC) {
      return traceToPath(denomTrace.value?.denom?.trace);
    } else {
      return "";
    }
  });

  const pathExtracted = computed(() => {
    if (isIBC) {
      return (
        traceToPath(denomTrace.value?.denom?.trace).match(/\d+/g)?.reverse() ?? ""
      );
    } else {
      return "";
    }
  });
  return { isIBC, denomTrace, normalized, path, pathExtracted };
};
export const useDenom = (denom: string) => {
  if (!useDenomInstances[denom]) {
    useDenomInstances[denom] = useDenomInstance(denom);
  }
  return useDenomInstances[denom];
};
