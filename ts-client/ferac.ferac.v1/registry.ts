import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgUpdateParams } from "./types/ferac/ferac/v1/tx";
import { MsgAnnounceValidatorExit } from "./types/ferac/ferac/v1/tx";
import { MsgSetValidatorQuality } from "./types/ferac/ferac/v1/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/ferac.ferac.v1.MsgUpdateParams", MsgUpdateParams],
    ["/ferac.ferac.v1.MsgAnnounceValidatorExit", MsgAnnounceValidatorExit],
    ["/ferac.ferac.v1.MsgSetValidatorQuality", MsgSetValidatorQuality],
    
];

export { msgTypes }