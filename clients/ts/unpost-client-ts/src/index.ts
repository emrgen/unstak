import {
    Configuration,
    DocumentServiceApi,
    GroupServiceApi,
    ProjectServiceApi,
    UserServiceApi,
    UserServiceApiCreateUserRequest,
    UserServiceApiDeleteUserRequest,
    UserServiceApiGetUserRequest,
    UserServiceApiListUsersRequest, UserServiceApiUpdateUserRequest
} from "@emrgen/unpost-client-gen";
import type {AxiosInstance} from "axios";
import axios from "axios";

class Config {
    token: () => string;
    basePath: string;
    axios?: AxiosInstance;
}

class documentClient {
    user: UserService;
    project: ProjectServiceApi;
    unpost: DocumentService;
    group: GroupServiceApi;

    constructor(config: Config) {
        const {basePath, globalAxios = axios} = config;
        const configuration = new Configuration({
            accessToken: config.token,
            basePath: config.basePath
        });

        this.user = new UserService(config);
        this.project = new ProjectServiceApi(configuration, basePath);
        this.unpost = new DocumentServiceApi(configuration, basePath);
        this.group = new GroupServiceApi(configuration, basePath);
    }
}

interface CreateUserArgs extends UserServiceApiCreateUserRequest {}
interface GetUserArgs extends UserServiceApiGetUserRequest {}
interface ListUserArgs extends UserServiceApiListUsersRequest {}
interface DeleteUserArgs extends UserServiceApiDeleteUserRequest {}
interface UpdateUserArgs extends UserServiceApiUpdateUserRequest {}

class UserService extends UserServiceApi {
    constructor(config: Config) {
        super(new Configuration({
            accessToken: config.token,
            basePath: config.basePath
        }));
    }

    async create(args: CreateUserArgs) {
        return await this.createUser(args)
    }

    async get(args: GetUserArgs) {
        return await this.getUser(args)
    }

    async list(args: ListUserArgs) {
        return await this.listUsers(args)
    }

    async update(args: UpdateUserArgs) {
        return await this.updateUser(args)
    }

    async delete(args: DeleteUserArgs) {
        return await this.deleteUser(args)
    }
}
