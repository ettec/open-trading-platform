package com.ettech.fixmarketsimulator.marketdataserver;

import io.grpc.*;

public class MdAuthInterceptor implements ServerInterceptor {

    static final String SubscriberIdKey = "subscriber_id";

    public static final Context.Key<String> SUBSCRIBER_ID
            = Context.key("subscriber_id"); // "identity" is just for debugging

    @Override
    public <ReqT, RespT> ServerCall.Listener<ReqT> interceptCall(
            ServerCall<ReqT, RespT> call,
            Metadata headers,
            ServerCallHandler<ReqT, RespT> next) {

        var key =  Metadata.Key.of(SubscriberIdKey, new Metadata.AsciiMarshaller<String>() {
            @Override
            public String toAsciiString(String s) {
                return s;
            }

            @Override
            public String parseAsciiString(String s) {
                return s;
            }
        });

        var value = headers.get(key);

        if( value == null) {
            call.close(Status.UNAUTHENTICATED.withDescription(" subscriber_id required "),
                    new Metadata());
        }

        Context context = Context.current().withValue(SUBSCRIBER_ID, value);
        return Contexts.interceptCall(context, call, headers, next);
    }
}
