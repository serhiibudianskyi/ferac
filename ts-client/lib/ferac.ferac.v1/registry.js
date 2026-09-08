import { MsgUpdateParams } from "./types/ferac/ferac/v1/tx";
import { MsgAnnounceValidatorExit } from "./types/ferac/ferac/v1/tx";
import { MsgSetValidatorQuality } from "./types/ferac/ferac/v1/tx";
const msgTypes = [
    ["/ferac.ferac.v1.MsgUpdateParams", MsgUpdateParams],
    ["/ferac.ferac.v1.MsgAnnounceValidatorExit", MsgAnnounceValidatorExit],
    ["/ferac.ferac.v1.MsgSetValidatorQuality", MsgSetValidatorQuality],
];
export { msgTypes };
